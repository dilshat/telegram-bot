function startsWith(str, substr) {
    return substr.length > 0 && substring(str, 0, substr.length) === substr;
}

function charAt(string, index) {
    var first = string.charCodeAt(index);
    var second;
    if (first >= 0xD800 && first <= 0xDBFF && string.length > index + 1) {
        second = string.charCodeAt(index + 1);
        if (second >= 0xDC00 && second <= 0xDFFF) {
            return string.substring(index, index + 2);
        }
    }
    return string[index];
}

function slice(string, start, end) {
    var accumulator = "";
    var character;
    var stringIndex = 0;
    var unicodeIndex = 0;
    var length = string.length;

    while (stringIndex < length) {
        character = charAt(string, stringIndex);
        if (unicodeIndex >= start && unicodeIndex < end) {
            accumulator += character;
        }
        stringIndex += character.length;
        unicodeIndex += 1;
    }
    return accumulator;
}

function toNumber(value, fallback) {
    if (value === undefined) {
        return fallback;
    } else {
        return Number(value);
    }
}

function substring(string, start, end) {
    var realStart = toNumber(start, 0);
    var realEnd = toNumber(end, string.length);
    if (realEnd == realStart) {
        return "";
    } else if (realEnd > realStart) {
        return slice(string, realStart, realEnd);
    } else {
        return slice(string, realEnd, realStart);
    }
}

function getUserID(message, callback){
    return callback? callback.From.ID : message.From.ID
}


function getUserName(message, callback) {
    var name

    if (callback) {
        name = callback.From.FirstName + ' ' + callback.From.LastName
        if (name.length == 1) {
            name = callback.From.Username
        }
    } else {
        name = message.From.FirstName + ' ' + message.From.LastName
        if (name.length == 1) {
            name = message.From.Username
        }
    }

    return name
}