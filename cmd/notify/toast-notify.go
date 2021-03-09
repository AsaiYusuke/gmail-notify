package main

import (
	"bytes"
	"fmt"
	"log"
)

type toastNotify struct {
	shell    *powerShell
	template *commandTemplate
}

type toastParameter struct {
	AppID       string
	Launch      string
	ImageSrc    string
	Title       string
	Message1    string
	Message2    string
	SilentSound bool
	Group       string
}

func (t *toastNotify) update(s *status) {
	if s.hasNewArrived {
		t.notify(s)
	} else if s.hasRemoveRequest {
		t.remove(s)
	}
}

func (t *toastNotify) notify(s *status) {
	var hasSilentSound bool
	if t.count(s) > 0 {
		t.remove(s)
		hasSilentSound = true
	}

	var buffer bytes.Buffer
	err := t.template.createTemplate.Execute(&buffer, t.createToastParameter(s, hasSilentSound))
	if err != nil {
		log.Fatalf("Unable to convert template: %v", err)
	}
	err = t.shell.syncExecute(buffer.String())
	if err != nil {
		log.Fatalf("Unable to push notify: %v", err)
	}
}

func (t *toastNotify) remove(s *status) {
	var buffer bytes.Buffer
	err := t.template.removeTemplate.Execute(&buffer, t.createToastParameter(s, false))
	if err != nil {
		log.Fatalf("Unable to convert template: %v", err)
	}
	err = t.shell.syncExecute(buffer.String())
	if err != nil {
		log.Fatalf("Unable to remove notify: %v", err)
	}
}

func (t *toastNotify) count(s *status) int64 {
	var buffer bytes.Buffer
	err := t.template.countTemplate.Execute(&buffer, t.createToastParameter(s, false))
	if err != nil {
		log.Fatalf("Unable to convert template: %v", err)
	}
	count, err := t.shell.syncExecuteCount(buffer.String())
	if err != nil {
		log.Fatalf("Unable to get history: %v", err)
	}
	return count
}

func (t *toastNotify) createLaunch(s *status) string {
	return toastBaseURL + s.getEmail()
}

func (t *toastNotify) createToastParameter(s *status, hasSilentSound bool) toastParameter {
	return toastParameter{
		AppID:       s.conf.AppID,
		Launch:      t.createLaunch(s),
		ImageSrc:    s.paths.getImageSrcFilePath(),
		Title:       s.conf.Title,
		Message1:    s.getEmail(),
		Message2:    fmt.Sprintf(`Unread: %d`, s.numOfUnreadIDs),
		SilentSound: hasSilentSound,
		Group:       s.getEmail(),
	}
}
