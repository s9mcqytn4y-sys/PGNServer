<?php

declare(strict_types=1);

namespace App\Support\Errors;

enum KodeKesalahanApi: string
{
    case SERVER_SEHAT = 'SERVER_SEHAT';
    case DATABASE_TIDAK_TERHUBUNG = 'DATABASE_TIDAK_TERHUBUNG';
    case VALIDASI_GAGAL = 'VALIDASI_GAGAL';
    case DATA_TIDAK_DITEMUKAN = 'DATA_TIDAK_DITEMUKAN';
    case KESALAHAN_SERVER = 'KESALAHAN_SERVER';
}
