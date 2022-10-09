package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

    "loan/x/loan/types"
)

func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
     // получаем контекст
	ctx := sdk.UnwrapSDKContext(goCtx)
	// получаем кредит с помощью сообщения в котором содержится id
    loan, found := k.GetLoan(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
    }
	// проверяем статус кредита
    if loan.State != "approved" {
        return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
    }
	// получаем кредитора с помощью филда в структуре
    lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	//получаем заёмщика
    borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	// проверяем отправителья сообщения только заёмщик этого кредита может его оплачивать 
    if msg.Creator != loan.Borrower {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot repay: not the borrower")
    }
	// получаем количество займа
    amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
    // получаем процент за кредит
	fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
	// получаем налог
    collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	//отправляем деньги с заёмщика к кредитора выплачивая кредит
    err := k.bankKeeper.SendCoins(ctx, borrower, lender, amount)
    if err != nil {
        return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
    }
	// отправляем процент кредита
    err = k.bankKeeper.SendCoins(ctx, borrower, lender, fee)
    if err != nil {
        return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
    }
	// отдаём заёмщику его залог
    err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)
    if err != nil {
        return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
    }
	// меняем статус кредита на оплаченный
    loan.State = "repayed"
	// меняем значение кредита в базе
    k.SetLoan(ctx, loan)

    return &types.MsgRepayLoanResponse{}, nil
}