# Kapan harus pake Session atau JWT

## Session

- Disaat aplikasi butuh menyimpan data session sehingg user bisa di logout melalui database
- Disaat kita ingin menyimpan data secara sementara di DB sehinggi bisa dilakukan analytic
- Disaat dibutuhkan data yang akses lebih cepat dari database

## JWT

- aplikasi stateless
- menghemat budget
- Implementasi lebih mudah dibandingkan Session jika aplikasi nya diatas load balancer
