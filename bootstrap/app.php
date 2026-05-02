<?php

declare(strict_types=1);

use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Symfony\Component\HttpKernel\Exception\MethodNotAllowedHttpException;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__.'/../routes/web.php',
        api: __DIR__.'/../routes/api.php',
        commands: __DIR__.'/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        $middleware->appendToGroup('api', [
            \Illuminate\Http\Middleware\HandleCors::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        $exceptions->shouldRenderJsonWhen(
            fn (Request $request, \Throwable $throwable): bool => $request->is('api/*') || $request->expectsJson(),
        );

        $exceptions->render(function (ValidationException $exception, Request $request) {
            if (! $request->is('api/*')) {
                return null;
            }

            $detailKesalahan = collect($exception->errors())
                ->flatMap(
                    fn (array $daftarPesan, string $field) => collect($daftarPesan)
                        ->map(fn (string $pesan) => [
                            'field' => $field,
                            'pesan' => $pesan,
                        ])
                )
                ->values()
                ->all();

            return ResponApi::gagal(
                pesan: 'Permintaan tidak dapat diproses',
                kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                detailKesalahan: $detailKesalahan,
                statusHttp: 422,
            );
        });

        $exceptions->render(function (NotFoundHttpException $exception, Request $request) {
            if (! $request->is('api/*')) {
                return null;
            }

            return ResponApi::gagal(
                pesan: 'Data yang diminta tidak ditemukan',
                kodeKesalahan: KodeKesalahanApi::DATA_TIDAK_DITEMUKAN,
                statusHttp: 404,
            );
        });

        $exceptions->render(function (MethodNotAllowedHttpException $exception, Request $request) {
            if (! $request->is('api/*')) {
                return null;
            }

            return ResponApi::gagal(
                pesan: 'Metode permintaan tidak didukung untuk endpoint ini',
                kodeKesalahan: KodeKesalahanApi::KESALAHAN_SERVER,
                statusHttp: 405,
            );
        });

        $exceptions->render(function (\Throwable $exception, Request $request) {
            if (! $request->is('api/*')) {
                return null;
            }

            return ResponApi::gagal(
                pesan: 'Terjadi kendala pada server. Silakan coba beberapa saat lagi.',
                kodeKesalahan: KodeKesalahanApi::KESALAHAN_SERVER,
                statusHttp: 500,
            );
        });
    })->create();
