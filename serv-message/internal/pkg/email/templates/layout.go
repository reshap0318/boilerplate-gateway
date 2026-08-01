package templates

import "fmt"

// Layout wraps HTML content in the shared page chrome (header, container
// styling, footer) used by every outgoing email — template-rendered or
// custom. Keeping this in one place means changing the look of emails never
// requires touching each template individually.
func Layout(appName, content string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f9fafb;
            margin: 0;
            padding: 40px 20px;
            color: #374151;
        }
        .container {
            max-width: 500px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            overflow: hidden;
        }
        .header {
            background-color: #4f46e5;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            color: #ffffff;
            margin: 0;
            font-size: 24px;
            font-weight: 600;
            letter-spacing: 0.5px;
        }
        .content {
            padding: 40px 30px;
        }
        .content h2 {
            margin-top: 0;
            color: #111827;
            font-size: 20px;
            font-weight: 600;
        }
        .content p {
            font-size: 15px;
            line-height: 1.6;
            margin-bottom: 24px;
        }
        .button-container {
            text-align: center;
            margin: 30px 0;
        }
        .button {
            display: inline-block;
            background-color: #4f46e5;
            color: #ffffff !important;
            text-decoration: none;
            padding: 14px 32px;
            border-radius: 8px;
            font-weight: 600;
            font-size: 16px;
        }
        .link-text {
            font-size: 13px;
            color: #6b7280;
            margin-bottom: 8px;
        }
        .token-box, .link {
            background-color: #f3f4f6;
            padding: 12px 16px;
            border-radius: 6px;
            font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
            font-size: 13px;
            color: #4b5563;
            word-break: break-all;
            margin-top: 0;
            margin-bottom: 24px;
        }
        .warning {
            background-color: #fefce8;
            border-left: 4px solid #eab308;
            padding: 16px;
            margin-top: 30px;
            border-radius: 0 6px 6px 0;
            font-size: 14px;
            color: #854d0e;
        }
        .info {
            background-color: #eff6ff;
            border-left: 4px solid #3b82f6;
            padding: 16px;
            margin-top: 30px;
            border-radius: 0 6px 6px 0;
            font-size: 14px;
            color: #1e3a8a;
        }
        .footer {
            text-align: center;
            padding: 24px;
            background-color: #f9fafb;
            color: #9ca3af;
            font-size: 13px;
            border-top: 1px solid #e5e7eb;
        }
        .footer p {
            margin: 4px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>%s</h1>
        </div>
        <div class="content">
            %s
        </div>
        <div class="footer">
            <p>&copy; %s. All rights reserved.</p>
            <p>Email ini dikirim secara otomatis, mohon tidak membalas.</p>
        </div>
    </div>
</body>
</html>
`, appName, content, appName)
}
