package service

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"gitlab.com/isteshkov/brute-force-protection/repositories"

	"google.golang.org/grpc"
)

var errorProducer = errors.NewProducer("GENERAL_ERROR")

func NewService(config *Config, subnetsRepo repositories.Subnets, l logging.Logger) *Service {
	service := &Service{
		Logger:            l,
		cfg:               config,
		subnetsRepository: subnetsRepo,
	}

	service.rpcListener = newRpc(service)
	service.profilingApi = buildProfilingApi()
	service.technicalApi = buildMetricsApi(service)
	return service
}

type Service struct {
	cfg *Config
	logging.Logger
	rpcListener  *grpc.Server
	technicalApi API
	profilingApi API

	subnetsRepository repositories.Subnets
}

func (s *Service) SetLogger(l logging.Logger) {
	s.Logger = l
}

func (s *Service) ListenAndServe() {
	lis, err := net.Listen("tcp", s.cfg.RpcPort)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	done := make(chan bool)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		sig := <-signals
		s.Logger.WithField("signal", sig).Info("get interrupting signal")
		done <- true
		return
	}()

	if s.cfg.ProfilingApiPort != "" {
		go func() {
			err = s.profilingApi.Run(s.cfg.ProfilingApiPort)
			if err != nil {
				s.Logger.Error(errorProducer.Wrap(err))
			}
			done <- true
			return
		}()

		s.Logger.WithField("pid", os.Getpid()).Info("Profiling API running on port %s", s.cfg.ProfilingApiPort)
	}

	if s.cfg.TechnicalApiPort != "" {
		go func() {
			err = s.technicalApi.Run(s.cfg.TechnicalApiPort)
			if err != nil {
				s.Logger.Error(errorProducer.Wrap(err))
			}
			done <- true
			return
		}()

		s.Logger.WithField("pid", os.Getpid()).Info("Technical API running on port %s", s.cfg.TechnicalApiPort)
	}

	go func() {
		err = s.rpcListener.Serve(lis)
		if err != nil {
			s.Logger.Error(errorProducer.Wrap(err))
		}
		done <- true
		return
	}()

	s.Logger.WithField("pid", os.Getpid()).Info("Rpc server is running on port %s", s.cfg.RpcPort)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				s.Logger.Info("Shutdown application")
				time.Sleep(time.Millisecond)
				return
			}
		}
	}()

	wg.Wait()
}
