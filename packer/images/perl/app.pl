use Mojolicious::Lite;

require "./fx.pl";

get '/' => sub {
  my $ctx = shift;
  my $res = fx($ctx);
  $ctx->render(json => $res);
};

post '/' => sub {
  my $ctx = shift;
  my $res = fx($ctx);
  $ctx->render(json => $res);
};

app->start;
