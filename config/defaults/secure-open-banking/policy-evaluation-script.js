function getIdmClientDetails() {
    return {
        "id": "{{ .Identity.IdmClientId }}",
        "secret": "{{ .Identity.IdmClientSecret }}",
        "endpoint": "http://am/am/oauth2/realms/root/realms/{{ .Identity.AmRealm }}/access_token",
        "scope": "fr:idm:*",
        "idmAdminUsername": "{{ .Ig.IgIdmUser }}",
        "idmAdminPassword": "{{ .Ig.IgIdmPassword }}"
    }
}

// Constants
const statusList = ["Authorised", "Consumed"];
const script_name = "policy_evaluation_script.js";
logger.message(script_name + ": starting")

const accountsAndTransactionsPermissions = [
    {name: "READACCOUNTSBASIC", property: {permission: "ReadAccountsBasic", requestType: "accounts"}},
    {name: "READACCOUNTSDETAIL", property: {permission: "ReadAccountsDetail", requestType: "accounts"}},
    {name: "READBALANCES", property: {permission: "ReadBalances", requestType: "balances"}},
    {name: "READBENEFICIARIESBASIC", property: {permission: "ReadBeneficiariesBasic", requestType: "beneficiaries"}},
    {name: "READBENEFICIARIESDETAIL", property: {permission: "ReadBeneficiariesDetail", requestType: "beneficiaries"}},
    {name: "READDIRECTDEBITS", property: {permission: "ReadDirectDebits", requestType: "direct-debits"}},
    {name: "READOFFERS", property: {permission: "ReadOffers", requestType: "offers"}},
    {name: "READPAN", property: {permission: "ReadPAN", requestType: ""}},
    {name: "READPARTY", property: {permission: "ReadParty", requestType: "party"}},
    {name: "READPARTIES", property: {permission: "ReadParty", requestType: "parties"}},
    {name: "READPARTYPSU", property: {permission: "ReadPartyPSU", requestType: "party"}},
    {name: "READPRODUCT", property: {permission: "ReadProducts", requestType: "product"}},
    {name: "READPRODUCTS", property: {permission: "ReadProducts", requestType: "products"}},
    {
        name: "READSCHEDULEDPAYMENTSBASIC",
        property: {permission: "ReadScheduledPaymentsBasic", requestType: "scheduled-payments"}
    },
    {
        name: "READSCHEDULEDPAYMENTSDETAIL",
        property: {permission: "ReadScheduledPaymentsDetail", requestType: "scheduled-payments"}
    },
    {
        name: "READSTANDINGORDERSBASIC",
        property: {permission: "ReadStandingOrdersBasic", requestType: "standing-orders"}
    },
    {
        name: "READSTANDINGORDERSDETAIL",
        property: {permission: "ReadStandingOrdersDetail", requestType: "standing-orders"}
    },
    {name: "READSTATEMENTSBASIC", property: {permission: "ReadStatementsBasic", requestType: "statements"}},
    {name: "READSTATEMENTSDETAIL", property: {permission: "ReadStatementsDetail", requestType: "statements"}},
    {name: "READTRANSACTIONSBASIC", property: {permission: "ReadTransactionsBasic", requestType: "transactions"}},
    {name: "READTRANSACTIONSCREDITS", property: {permission: "ReadTransactionsCredits", requestType: "transactions"}},
    {name: "READTRANSACTIONSDEBITS", property: {permission: "ReadTransactionsDebits", requestType: "transactions"}},
    {name: "READTRANSACTIONSDETAIL", property: {permission: "ReadTransactionsDetail", requestType: "transactions"}}
];

const paymentsIntents = [
    "domesticPaymentIntent", "domesticScheduledPaymentIntent", "domesticStandingOrderIntent",
    "internationalPaymentIntent", "internationalScheduledPaymentIntent", "internationalStandingOrderIntent"
];

function getPermissionAccountAndTransactions(name) {
    for (let i = 0; i < accountsAndTransactionsPermissions.length; i++) {
        if (accountsAndTransactionsPermissions[i].name === name || accountsAndTransactionsPermissions[i].property.permission == name) {
            return accountsAndTransactionsPermissions[i].property.permission
        }
    }
    return null
}

function getRequestTypeAccountAndTransactions(requestType) {
    for (let i = 0; i < accountsAndTransactionsPermissions.length; i++) {
        if (accountsAndTransactionsPermissions[i].property.requestType === requestType) return accountsAndTransactionsPermissions[i].property.requestType
    }
    return null
}

function dataAuthorisedAccountAndTransactions(permissions, requestType) {
    switch (requestType) {
        case "transactions":
            return ((permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSBASIC")) > -1 || permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSDETAIL")) > -1) && (permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSDEBITS")) > -1 || permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSCREDITS")) > -1))
        default:
            for (let i = 0; i < accountsAndTransactionsPermissions.length; i++) {
                if (requestType == accountsAndTransactionsPermissions[i].property.requestType && permissions.indexOf(accountsAndTransactionsPermissions[i].property.permission) > -1) return true
            }
    }
    return false
}

function dataAuthorised(permissions, requestType) {
    if (dataAuthorisedAccountAndTransactions(permissions, requestType))
        return true
    return false
}

function parseResourceUri() {
    const elements = resourceURI.split("/");
    return {
        "api": elements[6].indexOf("?") > -1 ? elements[6].substring(0, elements[6].indexOf("?")) : elements[6],
        "id": (elements.length > 7) ? elements[7] : null,
        "data": (elements.length > 8) ? elements[8] : null
    }
}

const p = "=";
const tab = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

function base64encode(ba) {
    let i;
    const s = [], l = ba.length;
    const rm = l % 3;
    const x = l - rm;
    for (i = 0; i < x;) {
        const t = ba[i++] << 16 | ba[i++] << 8 | ba[i++];
        s.push(tab.charAt((t >>> 18) & 0x3f));
        s.push(tab.charAt((t >>> 12) & 0x3f));
        s.push(tab.charAt((t >>> 6) & 0x3f));
        s.push(tab.charAt(t & 0x3f));
    }
    //	deal with trailers, based on patch from Peter Wood.
    switch (rm) {
        case 2: {
            const t = ba[i++] << 16 | ba[i++] << 8;
            s.push(tab.charAt((t >>> 18) & 0x3f));
            s.push(tab.charAt((t >>> 12) & 0x3f));
            s.push(tab.charAt((t >>> 6) & 0x3f));
            s.push(p);
            break;
        }
        case 1: {
            const t = ba[i++] << 16;
            s.push(tab.charAt((t >>> 18) & 0x3f));
            s.push(tab.charAt((t >>> 12) & 0x3f));
            s.push(p);
            s.push(p);
            break;
        }
    }
    return s.join("");	//	string
}

function base64decode(str) {
    const s = str.split(""), out = [];
    let l = s.length;
    while (s[--l] == p) {
    }	//	strip off trailing padding
    for (let i = 0; i < l;) {
        let t = tab.indexOf(s[i++]) << 18;
        if (i <= l) {
            t |= tab.indexOf(s[i++]) << 12
        }
        ;
        if (i <= l) {
            t |= tab.indexOf(s[i++]) << 6
        }
        ;
        if (i <= l) {
            t |= tab.indexOf(s[i++])
        }
        ;
        out.push((t >>> 16) & 0xff);
        out.push((t >>> 8) & 0xff);
        out.push(t & 0xff);
    }
    //	strip off any null bytes
    while (out[out.length - 1] == 0) {
        out.pop();
    }
    return out;	//	byte[]
}

function stringFromArray(data) {
    const count = data.length;
    let str = "";

    for (var index = 0; index < count; index += 1)
        str += String.fromCharCode(data[index]);

    return str;
}

function logResponse(callerMethod, response) {
    logger.message(script_name + ": [" + callerMethod + "] OB_Policy User REST Call. Status: " + response.getStatus() + ", Body: " + response.getEntity().getString());
}

function getIdmAccessToken() {

    const clientInfo = getIdmClientDetails();
    const request = new org.forgerock.http.protocol.Request();
    request.setUri(clientInfo.endpoint);
    request.setMethod("POST");
    request.getHeaders().add("Content-Type", "application/x-www-form-urlencoded");
    const formvars = "grant_type=password" +
        "&client_id=" + clientInfo.id +
        "&client_secret=" + clientInfo.secret +
        "&scope=" + clientInfo.scope +
        "&username=" + clientInfo.idmAdminUsername +
        "&password=" + clientInfo.idmAdminPassword;
    request.setEntity(formvars);

    const response = httpClient.send(request).get();


    logResponse("getIdmAccessToken", response);

    const oauth2response = JSON.parse(response.getEntity().getString());

    logger.message(script_name + ": Got access token " + oauth2response.access_token);
    return oauth2response.access_token
}

function findIntentType(api) {
    if (getRequestTypeAccountAndTransactions(api) != null) {
        return "accountAccessIntent"
    } else if (api === "domestic-payments" || api === "domestic-payment-consents") {
        return "domesticPaymentIntent"
    } else if (api === "domestic-scheduled-payments" || api === "domestic-scheduled-payment-consents") {
        return "domesticScheduledPaymentIntent"
    } else if (api === "domestic-standing-orders" || api === "domestic-standing-order-consents") {
        return "domesticStandingOrderIntent"
    } else if (api === "international-payments" || api === "international-payment-consents") {
        return "internationalPaymentIntent"
    } else if (api === "international-scheduled-payments" || api === "international-scheduled-payment-consents") {
        return "internationalScheduledPaymentIntent"
    } else if (api == "international-standing-orders" || api == "international-standing-order-consents") {
        return "internationalStandingOrderIntent"
    }
    return null
}

function getIntent(intentId, intentType) {
    const accessToken = getIdmAccessToken();
    const request = new org.forgerock.http.protocol.Request();
    const uri = "http://idm/openidm/managed/" + intentType + "/" + intentId + "?_fields=_id,_rev,Data,Risk,user/_id,accounts,apiClient/_id"
    logger.message(script_name + ": IDM fetch " + uri)

    request.setMethod('GET');
    request.setUri(uri)
    request.getHeaders().add("Authorization", "Bearer " + accessToken);

    const response = httpClient.send(request).get();
    logResponse("getIntent", response);

    return JSON.parse(response.getEntity().getString());
}

function deepCompare(arg1, arg2) {
    if (Object.prototype.toString.call(arg1) === Object.prototype.toString.call(arg2)) {
        if (Object.prototype.toString.call(arg1) === '[object Object]' || Object.prototype.toString.call(arg1) === '[object Array]') {
            if (Object.keys(arg1).length !== Object.keys(arg2).length) {
                return false;
            }
            return (Object.keys(arg1).every(function (key) {
                return deepCompare(arg1[key], arg2[key]);
            }));
        }
        return (arg1 === arg2);
    }
    return false;
}

function initiationMatch(initiationRequest, initiation) {
    const initiationRequestObj = JSON.parse(stringFromArray(base64decode(initiationRequest)))
    if (initiation.DebtorAccount && initiation.DebtorAccount.AccountId) {
        delete initiation.DebtorAccount.AccountId;
    }
    logger.message(script_name + ": initiationRequestObj " + JSON.stringify(initiationRequestObj))
    logger.message(script_name + ": initiation " + JSON.stringify(initiation))

    const match = deepCompare(initiationRequestObj, initiation);
    if (!match) {
        logger.warning(script_name + ": Mismatch between request [" + JSON.stringify(initiationRequestObj) + "] and consent [" + JSON.stringify(initiation) + "]");
    }

    return match
}

const intentId = environment.get("intent_id").iterator().next();
const apiRequest = parseResourceUri()
logger.message(script_name + ": req " + apiRequest.api + ":" + apiRequest.id + ":" + apiRequest.data);

const intentType = findIntentType(apiRequest.api)
const intent = getIntent(intentId, intentType);

const status = intent.Data.Status
const permissions = intent.Data.Permissions
const accounts = intent.accounts
// The responseAttributes expected always and array as value
const userResourceOwner = new Array(intent.user._id)

if (intentType === "accountAccessIntent") {
    logger.message(script_name + ": Account Access Intent");

    if (statusList.indexOf(status) == -1) {
        logger.message(script_name + "-[Account Access]: Rejecting request - status [" + status + "]")
        authorized = false
    } else if (apiRequest.id == null) {
        logger.message(script_name + "-[Account Access]: accounts " + accounts);
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = true
    } else if (apiRequest.data == null) {
        logger.message(script_name + "-[Account Access]: account info for " + apiRequest.id);
        // RS server expects granted accounts and permissions even though we're checking as well
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = (accounts.indexOf(apiRequest.id) > -1) &&
            dataAuthorised(permissions, apiRequest.api)
    } else {
        logger.message(script_name + "-[Account Access]: account request for " + apiRequest.id + ":" + apiRequest.data);
        // RS server expects granted accounts and permissions even though we're checking as well
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = (accounts.indexOf(apiRequest.id) > -1) &&
            dataAuthorised(permissions, apiRequest.data)
    }

} else if (paymentsIntents.indexOf(intentType) !== -1) {
    logger.message(script_name + ": Payments Intent");

    if (statusList.indexOf(status) == -1) {
        logger.message(script_name + "-[Payments]: Rejecting request - status [" + status + "]")
        authorized = false
    } else {
        responseAttributes.put("userResourceOwner", userResourceOwner);
        const requestMethod = environment.get("request_method").iterator().next();

        if (requestMethod != null) {
            switch (requestMethod) {
                case "POST":
                    logger.message(script_name + "-[Payments]: POST request");
                    const initiation = environment.get("initiation").iterator().next();
                    authorized = initiationMatch(initiation, intent.Data.Initiation)
                    break;
                case "GET":
                    logger.message(script_name + "-[Payments]: GET request");
                    authorized = true
                    break;
                default:
                    authorized = false
                    break;
            }
        } else {
            authorized = false
        }
    }
} else {
    authorized = false
}
logger.message(script_name + ": Policy evaluation result, authorized=" + authorized);
