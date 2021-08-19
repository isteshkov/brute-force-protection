package service

import (
	"net"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/domain/models"
	"gitlab.com/isteshkov/brute-force-protection/repositories"
)

func (s *Service) authAttempt(login, password, ip string) (allow bool, err error) {
	defer processError(&err)

	allow = false

	IPv4Addr, IPv4Subnet, err := net.ParseCIDR(ip)
	if err != nil {
		return
	}
	address, err := s.subnetsRepository.ByAddress(IPv4Subnet.String())
	switch {
	case err == nil:
		if !address.IsBlacklisted {
			allow = true
		}
		return
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
		err = nil
	default:
		return
	}

	if !s.rateLimiter.AttemptByLogin(login) {
		return
	}

	if !s.rateLimiter.AttemptByPassword(password) {
		return
	}

	if !s.rateLimiter.AttemptByIP(IPv4Addr.String()) {
		return
	}

	allow = true
	return
}

func (s *Service) cleanBucketByLogin(login string) (err error) {
	defer processError(&err)

	err = s.rateLimiter.CleanBucketByLogin(login)
	if err != nil {
		return
	}

	return
}

func (s *Service) cleanBucketByIP(ip string) (err error) {
	defer processError(&err)

	err = s.rateLimiter.CleanBucketByIP(ip)
	if err != nil {
		return
	}

	return
}

func (s *Service) addSubnetToWhitelist(ipAddress string) (err error) {
	defer processError(&err)

	var (
		existed models.Subnet
		subnet  models.Subnet
	)

	_, IPv4Subnet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return
	}

	existed, err = s.subnetsRepository.ByAddress(IPv4Subnet.String())
	switch {
	case err == nil:
		subnet.Version = existed.Version
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
		subnet.Address = IPv4Subnet.String()
	default:
		return
	}

	tx, err := s.subnetsRepository.Set(subnet, nil)
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	return
}

func (s *Service) addSubnetToBlacklist(ipAddress string) (err error) {
	defer processError(&err)

	var (
		existed models.Subnet
		subnet  models.Subnet
	)

	_, IPv4Subnet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return
	}

	existed, err = s.subnetsRepository.ByAddress(IPv4Subnet.String())
	switch {
	case err == nil:
		subnet.Version = existed.Version
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
		subnet.Address = IPv4Subnet.String()
	default:
		return
	}

	subnet.IsBlacklisted = true
	tx, err := s.subnetsRepository.Set(subnet, nil)
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	return
}

func (s *Service) removeSubnetFromList(ipAddress string) (err error) {
	defer processError(&err)

	_, IPv4Subnet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return
	}

	subnet, err := s.subnetsRepository.ByAddress(IPv4Subnet.String())
	if err != nil {
		return
	}

	tx, err := s.subnetsRepository.SetDeleted(subnet, time.Now(), nil)
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		if tx != nil {
			tx.MustRollBack(err.Error())
		}
		return
	}

	return
}
