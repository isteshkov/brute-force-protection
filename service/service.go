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

	service.rpcListener = newRPC(service)
	service.profilingAPI = buildProfilingAPI()
	service.technicalAPI = buildMetricsAPI(service)
	return service
}

type Service struct {
	cfg *Config
	logging.Logger
	rpcListener  *grpc.Server
	technicalAPI API
	profilingAPI API

	subnetsRepository repositories.Subnets
}

func (s *Service) SetLogger(l logging.Logger) {
	s.Logger = l
}

func (s *Service) ListenAndServe() {
	lis, err := net.Listen("tcp", s.cfg.RPCPort)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	done := make(chan bool)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-signals
		s.Logger.WithField("signal", sig).Info("get interrupting signal")
		done <- true
	}()

	if s.cfg.ProfilingAPIPort != "" {
		go func() {
			err = s.profilingAPI.Run(s.cfg.ProfilingAPIPort)
			if err != nil {
				s.Logger.Error(errorProducer.Wrap(err))
			}
			done <- true
		}()

		s.Logger.WithField("pid", os.Getpid()).Info("Profiling API running on port %s", s.cfg.ProfilingAPIPort)
	}

	if s.cfg.TechnicalAPIPort != "" {
		go func() {
			err = s.technicalAPI.Run(s.cfg.TechnicalAPIPort)
			if err != nil {
				s.Logger.Error(errorProducer.Wrap(err))
			}
			done <- true
		}()

		s.Logger.WithField("pid", os.Getpid()).Info("Technical API running on port %s", s.cfg.TechnicalAPIPort)
	}

	go func() {
		err = s.rpcListener.Serve(lis)
		if err != nil {
			s.Logger.Error(errorProducer.Wrap(err))
		}
		done <- true
	}()

	s.Logger.WithField("pid", os.Getpid()).Info("Rpc server is running on port %s", s.cfg.RPCPort)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range done {
			s.Logger.Info("Shutdown application")
			time.Sleep(time.Millisecond)
			return
		}
	}()

	wg.Wait()
}
