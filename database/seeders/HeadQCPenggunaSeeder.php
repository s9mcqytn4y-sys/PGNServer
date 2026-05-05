<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\User;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

final class HeadQCPenggunaSeeder extends Seeder
{
    use WithoutModelEvents;

    public function run(): void
    {
        /** @var array{email:string,password_default:string,nama_pengguna:string,peran:string} $konfigurasiHeadQC */
        $konfigurasiHeadQC = config('qcontrol.headqc');

        User::query()->updateOrCreate(
            [
                'email' => $konfigurasiHeadQC['email'],
            ],
            [
                'name' => $konfigurasiHeadQC['nama_pengguna'],
                'password' => Hash::make($konfigurasiHeadQC['password_default']),
                'peran' => 'HeadQC',
            ],
        );
    }
}
