package service

import (
	"fmt"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/domain/models"
	"gitlab.com/isteshkov/brute-force-protection/repositories"
)

func (s *Service) authAttempt(login, password, ip string) (err error) {
	defer processError(&err)

	// business logic here
	fmt.Println(login, password, ip)

	return
}

func (s *Service) cleanBucketByLogin(login string) (err error) {
	defer processError(&err)

	// business logic here
	_ = login

	return
}

func (s *Service) cleanBucketByIP(ip string) (err error) {
	defer processError(&err)

	// business logic here
	_ = ip

	return
}

func (s *Service) addSubnetToWhitelist(subnetAddress string) (err error) {
	defer processError(&err)

	var (
		existed models.Subnet
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
		existed models.Subnet
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
