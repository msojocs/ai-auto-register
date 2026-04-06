package mailprovider

import (
	"context"
	"regexp"
	"testing"
)

func TestLinshiyouxiangReceive(t *testing.T) {
	provider := NewLinshiyouxiang(map[string]string{
		"api_url": "https://deepmails.org",
	})
	ctx := context.Background()
	email, err := provider.GetEmail(ctx)
	if err != nil {
		t.Fatalf("GetEmail: %v", err)
	}
	t.Logf("Created email: %s", email)
	t.Logf("去 https://anonymousemail.me/ 发邮件给：%s", email.Email)
	t.Logf("内容：code: 123456")

	code, err := provider.WaitForCode(ctx, email, "code", 600)
	if err != nil {
		t.Fatalf("WaitForCode: %v", err)
	}
	t.Logf("Received messages: %v", code)
	if code != "123456" {
		t.Fatalf("Expected code 123456, got %s", code)
	}
}

func TestMailDetailSearch(t *testing.T) {
	data := `<div class="table-responsive msglist  pb-4">
                    <div class="mail-top clearfix">
                        <ul class="list-unstyled pull-left">
                            <li>
                                
                                    <a id="back-to-list" href="/zh">
                                    <span class="fa fa-chevron-left"></span><span data-lang-key="back">返回</span></a>
                            </li>
                        </ul>
                        <ul class="list-unstyled pull-right">
                            <li>
                                <a onclick="deleteMessage('0dce12941485e8202d997321256bf1e1')" href="javascript:void(0);">
                                    <i class="fa fa-close"></i><span data-lang-key="delete">删除</span></a>
                            </li>
                        </ul>
                    </div>
                    <div class="mail-info">
                        <h4 class="pt-2">你的 ChatGPT 代码为 705051</h4>
                        <ul class="list-unstyled">
                            <li>
                                <i class="fa fa-user-o"></i><span data-lang-key="from">发件人</span>: 
                                
                                    OpenAI &lt;<a href="/cdn-cgi/l/email-protection" class="__cf_email__" data-cfemail="80eff4f0c0f4edb1aeeff0e5eee1e9aee3efed">[email&#160;protected]</a>&gt;
                                
                            </li>
                            <li>
                                <i class="fa fa-clock-o"></i><span data-lang-key="time">时间</span>: <span id="time">1775451744</span>
                            </li>
                        </ul>
                    </div>
                    <div class="email-content">
                        <iframe id="email-iframe" 
                                sandbox="allow-same-origin allow-scripts allow-popups" 
                                style="width: 100%; border: none; overflow: hidden;min-height: 800px;"
                                scrolling="no">
                        </iframe>
                    </div>
            
                    
                </div>
                <hr>
                <div class="site-description">`
	detailReg := regexp.MustCompile(`<div class="table-responsive msglist  pb-4">([\s\S]+?)<div class="site-description">`)
	detailMatch := detailReg.FindSubmatch([]byte(data))
	if len(detailMatch) < 2 {
		t.Fatalf("failed to parse message detail from response")
	}
	detail := string(detailMatch[1])
	t.Logf("Parsed detail: %s", detail)
	if !regexp.MustCompile(`你的 ChatGPT 代码为 (\d{6})`).MatchString(detail) {
		t.Fatalf("failed to find code in detail")
	}

}
