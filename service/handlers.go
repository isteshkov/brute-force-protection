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
		if address != nil {
			if !address.IsBlacklisted {
				allow = true
			}
			return
		}
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
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

func (s *Service) addSubnetToWhitelist(subnetAddress string) (err error) {
	defer processError(&err)

	var (
		existed *models.Subnet
		subnet  models.Subnet
	)

	existed, err = s.subnetsRepository.ByAddress(subnetAddress)
	switch {
	case err == nil:
		subnet.Version = existed.Version
		fallthrough
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
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

func (s *Service) addSubnetToBlacklist(subnetAddress string) (err error) {
	defer processError(&err)

	var (
		existed *models.Subnet
		subnet  models.Subnet
	)

	existed, err = s.subnetsRepository.ByAddress(subnetAddress)
	switch {
	case err == nil:
		subnet.Version = existed.Version
		fallthrough
	case errors.IsProducedBy(err, repositories.ErrorProducerDoesNotExists):
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

func (s *Service) removeSubnetFromList(subnetAddress string) (err error) {
	defer processError(&err)

	subnet, err := s.subnetsRepository.ByAddress(subnetAddress)
	if err != nil {
		return
	}

	tx, err := s.subnetsRepository.SetDeleted(*subnet, time.Now(), nil)
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
