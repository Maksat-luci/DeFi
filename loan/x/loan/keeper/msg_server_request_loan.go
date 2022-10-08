package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"

    "loan/x/loan/types"
)

func (k msgServer) RequestLoan(goCtx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Создайте новый кредит со следующим пользовательским вводом
    var loan = types.Loan{
        Amount:     msg.Amount,
        Fee:        msg.Fee,
        Collateral: msg.Collateral,
        Deadline:   msg.Deadline,
        State:      "requested",
        Borrower:   msg.Creator,
    }


    // moduleAcc := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
    // //Получаем адресс заемщика из сообшения используя функцию AccAdressFromBech32
    borrower, _ := sdk.AccAddressFromBech32(msg.Creator)

    ////Получите залог как sdk.Coins
    collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
    if err != nil {
        panic(err)
    }

    // //Используйте учетную запись модуля в качестве условного счета
    sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, collateral)
    if sdkError != nil {
        return nil, sdkError
    }

    // //Добавить кредит в хранитель
    k.AppendLoan(
        ctx,
        loan,
    )

    return &types.MsgRequestLoanResponse{}, nil
}