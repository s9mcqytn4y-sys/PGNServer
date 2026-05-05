<?php

declare(strict_types=1);

return [
    'headqc' => [
        'email' => env('QCONTROL_HEADQC_EMAIL', 'headqc@pgn.local'),
        'password_default' => env('QCONTROL_HEADQC_PASSWORD', 'HeadQC@12345'),
        'nama_pengguna' => 'HeadQC',
        'peran' => 'HeadQC',
    ],
];
