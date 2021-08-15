package service

func (s *Service) authAttempt(login, password, ip string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) cleanBucketByLogin(login string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) cleanBucketByIP(ip string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) addSubnetToWhitelist(subnetAddress string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) addSubnetToBlacklist(subnetAddress string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) removeSubnetFromWhitelist(subnetAddress string) (err error) {
	defer processError(&err)

	// business logic here

	return
}

func (s *Service) removeSubnetFromBlacklist(subnetAddress string) (err error) {
	defer processError(&err)

	// business logic here

	return
}
