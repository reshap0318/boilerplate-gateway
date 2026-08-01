package templates

import "fmt"

// ResetPasswordContent renders the content section of the reset-password
// email. Wrap it with Layout before sending.
func ResetPasswordContent(token, resetURL string) string {
	return fmt.Sprintf(`
<h2>Permintaan Reset Password</h2>
<p>Halo,</p>
<p>Kami menerima permintaan untuk mengatur ulang kata sandi akun Anda. Anda dapat mereset kata sandi dengan mengklik tombol di bawah ini:</p>

<div class="button-container">
    <a href="%s" class="button">Reset Kata Sandi</a>
</div>

<p class="link-text">Atau salin dan tempel tautan berikut ke browser Anda:</p>
<p class="token-box" style="text-align: left; font-weight: normal; letter-spacing: normal;">%s</p>

<p class="link-text">Kode Token Manual:</p>
<p class="token-box" style="text-align: center; letter-spacing: 1px; font-weight: 600;">%s</p>

<div class="warning">
    <strong>Penting:</strong> Tautan ini hanya berlaku selama 1 jam. Jika Anda tidak meminta pengaturan ulang kata sandi, harap abaikan email ini atau segera hubungi layanan dukungan kami.
</div>
`, resetURL, resetURL, token)
}
