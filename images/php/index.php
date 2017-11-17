<?php
    include("fx.php");
    $data = file_get_contents("php://input");
    $res = json_decode($data,true);
    $v = Fx($res);
    echo $v;
