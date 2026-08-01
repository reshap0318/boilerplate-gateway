package templates

import "fmt"

// VerifyEmailContent renders the content section of the email-verification
// email. Wrap it with Layout before sending.
func VerifyEmailContent(verifyURL, appName string) string {
	return fmt.Sprintf(`
<h2>Verifikasi Email Anda</h2>
<p>Halo,</p>
<p>Terima kasih telah mendaftar di <strong>%s</strong>! Untuk menyelesaikan proses registrasi dan mengamankan akun Anda, silakan verifikasi alamat email Anda dengan mengklik tombol di bawah ini:</p>

<div class="button-container">
    <a href="%s" class="button">Verifikasi Email Sekarang</a>
</div>

<p class="link-text">Atau salin dan tempel tautan berikut ke browser Anda:</p>
<p class="link">%s</p>

<div class="info">
    <strong>Info:</strong> Tautan verifikasi ini hanya berlaku selama 24 jam. Jika Anda merasa tidak mendaftar di %s, Anda dapat mengabaikan email ini dengan aman.
</div>
`, appName, verifyURL, verifyURL, appName)
}
