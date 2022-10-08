package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//получаем кредит с хранилища при помощи Айди в сообщении
	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
	}

	//TODO: по какой-то причине ошибка не выводится на терминал
	//проверяем его статус
	if loan.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}
	//lender получаем кредитора из сообщения
	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	// получаем адресс человека который попросил кредит
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	//получаем деньги необходимык заёмщику
	amount, err := sdk.ParseCoinsNormalized(loan.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot parse coins in loan amount")
	}
	// отправляем деньги с кредитора к заёмщику
	k.bankKeeper.SendCoins(ctx, lender, borrower, amount)
	// меняем статус на одобренный и присваеваем нового создателя, хз почему кредитора но ладно
	loan.Lender = msg.Creator
	loan.State = "approved"
	// изменяем кредит в хранилище
	k.SetLoan(ctx, loan)

	return &types.MsgApproveLoanResponse{}, nil
}
