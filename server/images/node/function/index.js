module.exports = (req) => {
    return req.querystring ? req.querystring.length : 0
}
