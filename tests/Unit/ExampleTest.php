<?php

declare(strict_types=1);

test('enum kode kesalahan api memiliki nilai stabil', function () {
    expect(App\Support\Errors\KodeKesalahanApi::KESALAHAN_SERVER->value)
        ->toBe('KESALAHAN_SERVER');
});
