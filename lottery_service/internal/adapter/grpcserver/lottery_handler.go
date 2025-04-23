package grpcserver

import (
	"context"

	"github.com/CharLottery/lottery_service/internal/model"
	"github.com/CharLottery/lottery_service/internal/usecase"
	lotterypb "github.com/CharLottery/proto/lotterypb"
)

type LotteryHandler struct {
	usecase *usecase.LotteryUsecase
	lotterypb.UnimplementedLotteryServiceServer
}

func NewLotteryHandler(uc *usecase.LotteryUsecase) *LotteryHandler {
	return &LotteryHandler{usecase: uc}
}

func (h *LotteryHandler) CreateLottery(ctx context.Context, req *lotterypb.CreateLotteryRequest) (*lotterypb.CreateLotteryResponse, error) {
	lottery := &model.Lottery{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Prize:       req.GetPrize(),
		Status:      "open",
	}
	if err := h.usecase.CreateLottery(lottery); err != nil {
		return nil, err
	}
	return &lotterypb.CreateLotteryResponse{
		Lottery: &lotterypb.Lottery{
			Id:          lottery.ID,
			Title:       lottery.Title,
			Description: lottery.Description,
			Prize:       lottery.Prize,
			Status:      lottery.Status,
		},
	}, nil
}

func (h *LotteryHandler) GetLottery(ctx context.Context, req *lotterypb.GetLotteryRequest) (*lotterypb.GetLotteryResponse, error) {
	lottery, err := h.usecase.GetLottery(req.GetId())
	if err != nil {
		return nil, err
	}
	return &lotterypb.GetLotteryResponse{
		Lottery: &lotterypb.Lottery{
			Id:          lottery.ID,
			Title:       lottery.Title,
			Description: lottery.Description,
			Prize:       lottery.Prize,
			Status:      lottery.Status,
		},
	}, nil
}

func (h *LotteryHandler) ListLotteries(ctx context.Context, req *lotterypb.ListLotteriesRequest) (*lotterypb.ListLotteriesResponse, error) {
	lotteries, err := h.usecase.ListLotteries()
	if err != nil {
		return nil, err
	}
	var result []*lotterypb.Lottery
	for _, l := range lotteries {
		result = append(result, &lotterypb.Lottery{
			Id:          l.ID,
			Title:       l.Title,
			Description: l.Description,
			Prize:       l.Prize,
			Status:      l.Status,
		})
	}
	return &lotterypb.ListLotteriesResponse{Lotteries: result}, nil
}
