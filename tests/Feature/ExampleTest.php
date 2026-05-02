<?php

declare(strict_types=1);

test('root web mengembalikan json sederhana', function () {
    $this->getJson('/')
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'PGNServer REST API siap digunakan',
        ])
        ->assertJsonPath('data.namaAplikasi', 'PGNServer')
        ->assertJsonPath('data.versiApi', 'v1');
});
