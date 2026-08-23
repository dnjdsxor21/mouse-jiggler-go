//go:build darwin && cgo

// Package platform exposes the macOS input APIs required by the CLI.
package platform

/*
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
#include <unistd.h>

static int mouseJigglerTrusted(void) {
	return AXIsProcessTrusted() ? 1 : 0;
}

static int mouseJigglerJiggle(void) {
	CGEventRef current = CGEventCreate(NULL);
	if (current == NULL) {
		return 0;
	}

	CGPoint location = CGEventGetLocation(current);
	CFRelease(current);

	CGEventRef outbound = CGEventCreateMouseEvent(
		NULL,
		kCGEventMouseMoved,
		CGPointMake(location.x + 1.0, location.y),
		kCGMouseButtonLeft
	);
	if (outbound == NULL) {
		return 0;
	}
	CGEventPost(kCGHIDEventTap, outbound);
	CFRelease(outbound);

	usleep(20000);

	CGEventRef restoration = CGEventCreateMouseEvent(
		NULL,
		kCGEventMouseMoved,
		location,
		kCGMouseButtonLeft
	);
	if (restoration == NULL) {
		return 0;
	}
	CGEventPost(kCGHIDEventTap, restoration);
	CFRelease(restoration);

	return 1;
}
*/
import "C"

import "errors"

// Mouse reports permission and sends non-clicking mouse-move events.
type Mouse struct{}

// Trusted reports whether the current process has macOS Accessibility access.
func (Mouse) Trusted() bool {
	return C.mouseJigglerTrusted() == 1
}

// Jiggle moves one point and restores the original pointer position.
func (Mouse) Jiggle() error {
	if C.mouseJigglerJiggle() != 1 {
		return errors.New("could not create a macOS mouse event")
	}
	return nil
}
