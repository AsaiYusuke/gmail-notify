package main

type status struct {
	paths            *paths
	conf             *config
	subDir           string
	idsFile          *jsonAccess
	gmail            *gmailAPI
	toast            *toastNotify
	email            string
	prevUnreadIDs    []string
	unreadTimestamp  int64
	hasNewArrived    bool
	hasRemoveRequest bool
	numOfUnreadIDs   int
}

func (s *status) getEmail() string {
	if s.email == "" {
		s.email = s.gmail.getEmail()
	}
	return s.email
}

func (s *status) getPrevMessageIDs() []string {
	if s.unreadTimestamp != s.idsFile.getTimestamp() {
		s.unreadTimestamp = s.idsFile.readIDs(&s.prevUnreadIDs)
	}
	return s.prevUnreadIDs
}

func (s *status) setPrevMessageIDs(prevMessageIDs []string) {
	s.prevUnreadIDs = append([]string{}, prevMessageIDs...)
	s.unreadTimestamp = s.idsFile.write(s.prevUnreadIDs)
}

func (s *status) update() {
	latestUnreadIDs := s.gmail.getUnreadMessageIDs()

	prevUnreadIDs := s.getPrevMessageIDs()

	s.hasNewArrived = false
	for _, unreadMessageID := range latestUnreadIDs {
		var isExist bool
		for _, prevMessageID := range prevUnreadIDs {
			if unreadMessageID == prevMessageID {
				isExist = true
				break
			}
		}
		if !isExist {
			s.hasNewArrived = true
			break
		}
	}

	s.numOfUnreadIDs = len(latestUnreadIDs)
	s.hasRemoveRequest = len(latestUnreadIDs) < len(prevUnreadIDs)

	if s.hasNewArrived || s.hasRemoveRequest {
		s.setPrevMessageIDs(latestUnreadIDs)
	}

	s.toast.update(s)
}
