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
        User::query()->updateOrCreate(
            [
                'email' => 'headqc@pgn.local',
            ],
            [
                'name' => 'HeadQC',
                'password' => Hash::make('HeadQC@12345'),
                'peran' => 'HeadQC',
            ],
        );
    }
}
