<?php

declare(strict_types=1);

namespace App\Console\Commands;

use App\Models\User;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Hash;

final class PastikanHeadQCPenggunaCommand extends Command
{
    protected $signature = 'qcontrol:pastikan-headqc {--tanpa-reset-password : Jangan ubah password jika pengguna HeadQC sudah ada}';

    protected $description = 'Membuat atau memperbarui pengguna HeadQC default untuk runtime lokal QControl.';

    public function handle(): int
    {
        /** @var array{email:string,password_default:string,nama_pengguna:string,peran:string} $konfigurasiHeadQC */
        $konfigurasiHeadQC = config('qcontrol.headqc');

        $email = $konfigurasiHeadQC['email'];
        $namaPengguna = $konfigurasiHeadQC['nama_pengguna'];
        $peran = $konfigurasiHeadQC['peran'];
        $kataSandiBawaan = $konfigurasiHeadQC['password_default'];
        $tanpaResetPassword = (bool) $this->option('tanpa-reset-password');

        $pengguna = User::query()->where('email', $email)->first();
        $passwordDiaturUlang = false;

        if ($pengguna === null) {
            $pengguna = User::query()->create([
                'name' => $namaPengguna,
                'email' => $email,
                'password' => Hash::make($kataSandiBawaan),
                'peran' => $peran,
            ]);

            $passwordDiaturUlang = true;
        } else {
            $atributPerubahan = [
                'name' => $namaPengguna,
                'peran' => $peran,
            ];

            if (! $tanpaResetPassword) {
                $atributPerubahan['password'] = Hash::make($kataSandiBawaan);
                $passwordDiaturUlang = true;
            }

            $pengguna->forceFill($atributPerubahan)->save();
        }

        $this->info('Bootstrap HeadQC selesai.');
        $this->line('Email: '.$pengguna->email);
        $this->line('Peran: '.$pengguna->peran);
        $this->line('Password default diatur ulang: '.($passwordDiaturUlang ? 'ya' : 'tidak'));

        if (app()->environment(['local', 'testing'])) {
            $this->line('Sumber password default: konfigurasi lokal aktif.');
        }

        return self::SUCCESS;
    }
}
