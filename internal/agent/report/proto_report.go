// Модуль отправки данных на сервер
package report

import (
	"context"
	"sync"
	"time"

	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/agent/metrics"
	pb "github.com/romanmendelproject/go-yandex-metrics/proto"
	"github.com/romanmendelproject/go-yandex-metrics/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ReportBatchMetric отправка нескольких метрик в одном пакете в формате JSON
func ReportBatchMetricProto(ctx context.Context, cfg *config.ClientFlags, wg *sync.WaitGroup, metricsChannel <-chan *[]metrics.Metric) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Info("Closing report program")
			return
		case data := <-metricsChannel:
			if err := sendMetricProto(ctx, cfg.FlagReqAddr, *data); err != nil {
				log.Error(err)
			}

		}
	}
}

func updateMS(ctx context.Context, c pb.MetricsClient, in *pb.UpdateBatchRequest) error {
	_, err := c.UpdateBatch(ctx, in)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}
	log.Info("gRPC agent ", "updateMS")
	return nil

}

func sendMetricProto(ctx context.Context, addr string, metrics []metrics.Metric) error {

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("gRPC agent sendRequestMetricGRPC: did not connect: grpc.Dial, ", "about ERR"+err.Error())
		return err
	}
	defer conn.Close()
	c := pb.NewMetricsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	ms := []*pb.Metric{}

	for _, m := range metrics {
		ms = append(ms, &pb.Metric{
			ID:    m.ID,
			MType: string(m.MType),
			Delta: utils.UnPointer(m.Delta),
			Value: utils.UnPointer(m.Value),
		})
	}
	mss := &pb.UpdateBatchRequest{Metric: ms}
	updateMS(ctx, c, mss)

	return nil
}
