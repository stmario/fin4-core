module Actions.Update exposing (update)

import Actions.Command
    exposing
        ( addRewardsCmd
        , approveProposalCmd
        , likeCmd
        , loadActionsCmd
        , submitProposalCmd
        )
import Actions.Model exposing (Model)
import Actions.Msg exposing (Msg(..))
import Actions.Ports exposing (fileSelected)
import Debug
import Dict exposing (Dict)
import Main.Context exposing (Context)
import Material
import Model.Actions exposing (actionsDecoder)


update : Context -> Msg -> Model -> ( Model, Cmd Msg )
update ctx msg model =
    case msg of
        OnLoadActionsResponse resp ->
            case resp of
                Ok actions ->
                    { model | actions = Just actions, error = Nothing } ! []

                Err error ->
                    Debug.log (toString error)
                        { model
                            | actions = Nothing
                            , error = Just "Error loading your actions, please try again!"
                        }
                        ! []

        AddRewards actionId amount ->
            model ! [ addRewardsCmd ctx model actionId amount ]

        OnAddRewardsResponse resp ->
            model ! [ loadActionsCmd ctx 1 ]

        TickerTimout ->
            model ! [ loadActionsCmd ctx 1 ]

        SubmitProposal actionId proposal ->
            let
                props =
                    Dict.remove actionId model.claims

                imgs =
                    Dict.remove actionId model.claimImages
            in
            { model | claims = props, claimImages = imgs } ! [ submitProposalCmd ctx model actionId proposal ]

        ApproveClaim claimId ->
            model ! [ approveProposalCmd ctx model claimId ]

        ApproveClaimResponse resp ->
            model ! [ loadActionsCmd ctx 1 ]

        SetClaim tokenId proposal ->
            let
                props =
                    Dict.insert tokenId proposal model.claims
            in
            { model | claims = props } ! []

        OnDoLikeResponse resp ->
            model ! []

        DoLike tokenId state ->
            model ! [ likeCmd ctx tokenId state ]

        ImageSelected tokenId ->
            model ! [ fileSelected <| "imgclaim-" ++ toString tokenId ]

        ImageRead data ->
            let
                newImage =
                    { contents = data.contents
                    , filename = data.filename
                    }

                imgs =
                    Dict.insert data.tokenId newImage model.claimImages
            in
            { model
                | claimImages = imgs
            }
                ! [ Cmd.none ]

        Mdl msg_ ->
            Material.update Mdl msg_ model
