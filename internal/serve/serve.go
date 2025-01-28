package serve

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Serve struct {
	httpSrv *http.Server
}

func NewServe(port string, handler http.Handler) *Serve {
	return &Serve{
		httpSrv: &http.Server{
			Addr:         port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// RunServe запускает HTTP-сервер с поддержкой graceful shutdown.
// port - порт, на котором будет запущен сервер (например, ":8080").
// handler - обработчик запросов (например, *chi.Mux).
// Возвращает ошибку, если сервер не удалось запустить или завершить.
func (s *Serve) RunServe() error {
	const op = "serve.RunServe"
	// Настройка HTTP-сервера

	// Запуск сервера в горутине
	serverErr := make(chan error, 1)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("%s: ошибка запуска сервера: %v", op, err)
		}
		log.Println(fmt.Sprintf("Сервер доступен по адресу: s%s", s.httpSrv.Addr), "op", op)
	}()

	// Ожидание graceful shutdowns
	if err := s.waitForShutdown(); err != nil {
		return fmt.Errorf("%s ошибка при завершении работы сервера: %v", op, err)
	}

	log.Println("Сервер успешно завершил работу", "op", op)
	return nil
}

// waitForShutdown ожидает сигналов завершения и выполняет graceful shutdown сервера.
// Возвращает ошибку, если shutdown не удался.
func (s *Serve) waitForShutdown() error {
	const op = "serve.waitForShutdown"
	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала завершения
	<-shutdown
	log.Println("Завершение работы сервера...", "op", op)

	// Graceful shutdown для HTTP-сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: ошибка при завершении работы сервера: %v", op, err)
	}

	return nil
}
