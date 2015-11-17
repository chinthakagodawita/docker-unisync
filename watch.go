package main

import (
	"github.com/mattes/fsevents"
	"sort"
	"strings"
	"time"
)

// A list of all fsevents flags mapped to readable descriptions.
// @see https://developer.apple.com/library/mac/documentation/Darwin/Reference/FSEvents_Ref/index.html#//apple_ref/doc/c_ref/kFSEventStreamEventFlagNone
var noteDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs: "MustScanSubdirs",
	fsevents.UserDropped:     "UserDropped",
	fsevents.KernelDropped:   "KernelDropped",
	fsevents.EventIDsWrapped: "EventIDsWrapped",
	fsevents.HistoryDone:     "HistoryDone",
	fsevents.RootChanged:     "RootChanged",
	fsevents.Mount:           "Mount",
	fsevents.Unmount:         "Unmount",

	fsevents.ItemCreated:       "Created",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

const WATCH_POLL_INTERVAL = 5

func Watch(path string, ignored []string, poll bool, callback func(id uint64, path string, flags []string)) {
	dev, _ := fsevents.DeviceForPath(path)
	fsevents.EventIDForDeviceBeforeTime(dev, time.Now())

	es := &fsevents.EventStream{
		Paths:   []string{path},
		Latency: 50 * time.Millisecond,
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}
	es.Start()
	ec := es.Events

	poller := time.NewTicker(WATCH_POLL_INTERVAL * time.Second)

	// Disable the poller if polling isn't requested.
	if !poll {
		poller.Stop()
	}

	for {
		select {
		case event := <-ec:
			if FileIsIgnored(event.Path, ignored) {
				continue
			}

			flags := make([]string, 0)
			for bit, description := range noteDescription {
				if event.Flags&bit == bit {
					// LogDebug(description)
					flags = append(flags, description)
				}
			}
			sort.Sort(sort.StringSlice(flags))
			go callback(event.ID, event.Path, flags)
			es.Flush(false)
		case <-poller.C:
			go callback(0, "", []string{})
		}
	}
}

// Heavily inspired by github.com/ryanuber/go-glob.
func FileIsIgnored(file string, ignored []string) bool {
	for _, pattern := range ignored {
		if pattern == "" {
			continue
		}

		if pattern == IGNORE_WILDCARD {
			return true
		}

		parts := strings.Split(pattern, IGNORE_WILDCARD)

		// If no wildcards exist in the pattern, check for exact equality.
		if len(parts) == 1 {
			if file == pattern {
				return true
			}
		}

		leadingWildcard := strings.HasPrefix(pattern, IGNORE_WILDCARD)
		trailingWildcard := strings.HasSuffix(pattern, IGNORE_WILDCARD)
		end := len(parts) - 1

		// Check each part of the pattern for a match.
		for i, part := range parts {
			switch i {
			case 0:
				if leadingWildcard {
					// Skip if we're checking everything but the start of the string.
					continue
				} else if strings.HasPrefix(file, part) {
					return true
				}
			case end:
				if len(file) > 0 && !trailingWildcard && strings.HasSuffix(file, part) {
					// If the filepath isn't empty, we're not checking for the end of a
					// string and we have a suffix match, return true.
					return true
				}
			default:
				if strings.Contains(file, part) {
					return true
				}
			}
		}
	}

	return false
}
