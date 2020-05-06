// 错误处理
function handleError(xhr) {
    let code = 0;
    try {
        code = xhr.responseJSON.Code
    } catch (e) {
        typeof ChaosFunctions === "object" && ChaosFunctions.Logger({ Type: 'error', Title: "Handle Error：", Info: e });
    }

    if (code === -1) {
        alert(xhr.responseJSON.Error)
        return
    }

    let tips = "未知错误";
    if (typeof ChaosCodes === "object" && ChaosCodes[code]) {
        tips = ChaosCodes[code]
    }

    switch (xhr.status) {
        case 400:
            alert("[ 请求错误 ] " + tips);
            break;
        case 401:
            alert("[ 未授权 ] " + tips);
            break;
        case 403:
            alert("[ 拒绝访问 ] " + tips);
            break;
        case 404:
            alert("[ 请求出错 ] " + tips);
            break;
        case 408:
            alert("[ 请求超时 ] " + tips);
            break;
        case 500:
            alert("[ 服务器内部错误 ] " + tips);
            break;
        case 501:
            alert("[ 服务未实现 ] " + tips);
            break;
        case 502:
            alert("[ 网关错误 ] " + tips);
            break;
        case 503:
            alert("[ 服务不可用 ] " + tips);
            break;
        case 504:
            alert("[ 网关超时 ] " + tips);
            break;
        case 505:
            alert("[ HTTP版本不受支持 ] " + tips);
            break;
        default: break
    }
}
