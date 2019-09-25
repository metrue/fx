#
# ctx = {
#   :request => request,
#   :response => response,
#   :status => status,
#   :headers => headers,
# }
def fx(ctx)
  ctx[:response].body = "hello world"
end
