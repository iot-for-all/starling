export function getErrorMessage(ex, defaultMessage) {
    if (ex) {
        if (ex.response && ex.response.data) {
            return ex.response.data;
        } else if (ex.message) {
            return ex.message;
        }
    }

    return defaultMessage;
}

export function formatNumber(num) {
    if (!num) {
        return "0";
    }
    return num.toString().replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1,')
}

export function formatCount(collection, name) {
    if (collection) {
        if (collection.length === 0) {
            return " no " + name + "s found";
        } else if (collection.length === 1) {
            return " 1 " + name + " found";
        }
        return " " + collection.length + " " + name + "s found";
    }
    return " no " + name + "s found";
}