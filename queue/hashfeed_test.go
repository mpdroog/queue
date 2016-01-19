package queue

import (
	"crypto/md5"
	"testing"
)

// https://github.com/jpmens/diablo/blob/master/lib/hashfeed.c#L223
func TestMatch(t *testing.T) {
	feed := "1-120/360:8"
	valid := []string{
		"<part29of143.RndMw4FFWQ9TdIWjYDmt@camelsystem-powerpost.local>",
		"<part107of143.qXxdN72wwtbVqGCTy4c0@camelsystem-powerpost.local>",
		"<Part211of211.1E89F6A757E540EA91B7A0B4950B183D@1452090703.local>",
		"<part100of143.sz1ngV671ts7LFN9BgOS@camelsystem-powerpost.local>",
		"<part80of143.8&j$YJ9Q$amj0puxjX3A@camelsystem-powerpost.local>",
		"<part103of143.qXxdN72wwtbVqGCTy4c0@camelsystem-powerpost.local>",
		"<part35of143.RndMw4FFWQ9TdIWjYDmt@camelsystem-powerpost.local>",
		"<part32of143.RndMw4FFWQ9TdIWjYDmt@camelsystem-powerpost.local>",
		"<1452088368.57488.1@usnews.blocknews.net>",
		"<part105of143.8zTGFGGlYt5qGOZkDuW9@camelsystem-powerpost.local>",
		"<part78of143.MA3lszmuS6Rg3ykqmZY2@camelsystem-powerpost.local>",
		"<n6j67r$mrr$1@ns2.nl2k.ab.ca>",
		"<part81of143.8&j$YJ9Q$amj0puxjX3A@camelsystem-powerpost.local>",
		"<part116of143.8zTGFGGlYt5qGOZkDuW9@camelsystem-powerpost.local>",
		"<gISdnZC6e7ingRDLnZ2dnUU78eOdnZ2d@giganews.com>",
		"<Part196of211.1569634777FC4766869D939CDF824583@1452090703.local>",
		"<LQXN9Vnbv4tJBZ12gjCB_91o168@JBinUp.local>",
		"<part104of143.8zTGFGGlYt5qGOZkDuW9@camelsystem-powerpost.local>",
		"<Part163of261.583B442ADD4E407893CE61BA230960EE@1452087846.local>",
		"<568d1aa0$1$8947$edfe2c6d@reader.snellerdownloaden.com>",
		"<O70K3htnDXFedFGOOjh6@JBinUp.local>",
		"<qk1Kj0zXMjpPtxRuNfK@JBinUp.local>",
		"<l8oTcR0ZWezHoDLH6P8j@JBinUp.local>",
		"<WmAV9GZrhlzdNq4oQ1xB_115o168@JBinUp.local>",
		"<moyjlKYGXKKUlDC4QxF6_88o168@JBinUp.local>",
		"<Part151of211.EED3B94EBA254587916A850030D3F759@1452090703.local>",
		"<G8mToWiM7hAGOlHLxjt_115o168@JBinUp.local>",
	}
	invalid := []string{
		"<part6of143.htzsE$6rldAt7WuGTpr9@camelsystem-powerpost.local>",
		"<part60of143.RndMw4FFWQ9TdIWjYDmt@camelsystem-powerpost.local>",
		"<part2of143.htzsE$6rldAt7WuGTpr9@camelsystem-powerpost.local>",
		"<part5of143.htzsE$6rldAt7WuGTpr9@camelsystem-powerpost.local>",
		"<part116of143.MA3lszmuS6Rg3ykqmZY2@camelsystem-powerpost.local>",
		"<part21of143.JCJAPg5O$$aZxFMzFpEZ@camelsystem-powerpost.local>",
		"<part24of143.JCJAPg5O$$aZxFMzFpEZ@camelsystem-powerpost.local>",
		"<part141of143.qXxdN72wwtbVqGCTy4c0@camelsystem-powerpost.local>",
		"<Part33of137.D3EEDFF4539F491CB513980D156BF6F9@1452082593.local>",
		"<part117of143.MA3lszmuS6Rg3ykqmZY2@camelsystem-powerpost.local>",
		"<part37of143.3cyACN2UfykYsavFok&c@camelsystem-powerpost.local>",
		"<part35of143.JCJAPg5O$$aZxFMzFpEZ@camelsystem-powerpost.local>",
		"<part125of143.sz1ngV671ts7LFN9BgOS@camelsystem-powerpost.local>",
		"<twJloG2l1pN5Ok5dMWl0qcpSCMEQA7g2.820-820@JBinDown.local>",
		"<part126of143.sz1ngV671ts7LFN9BgOS@camelsystem-powerpost.local>",
		"<part121of143.8&j$YJ9Q$amj0puxjX3A@camelsystem-powerpost.local>",
		"<1452088239.92474.17@reader.easyusenet.nl>",
		"<FtyNAoeMxRnXHfFxxOvQl8Bj6rC2i1CB4.94-273@JBinDown.local>",
		"<part25of143.htzsE$6rldAt7WuGTpr9@camelsystem-powerpost.local>",
		"<part116of143.8&j$YJ9Q$amj0puxjX3A@camelsystem-powerpost.local>",
		"<part68of137.iuzjOe3evX4zLo1dTqAY@camelsystem-powerpost.local>",
		"<pelUI2u9nTRG7wcfFPR1_123o168@JBinUp.local>",
		"<Part46of211.720CD1911FD441CE96022C45FB74A8A4@1452090703.local>",
		"<Part238of261.AB0AF33AFB784BFC89F21D3A352F69CC@1452087846.local>",
		"<Part175of266.DD74BDCF8F4247EEB35AB1D8258B4E7F@1452087949.local>",
		"<YXY3SvIOI67dqnHnlSMV_83o100@JBinUp.local>",
		"<Part107of137.D71179EAFEAD41F498B4AD5C621AB8DE@1452109356.local>",
	}

	match, e := parseFeed(feed)
	if e != nil {
		t.Fatal(e)
	}

	for _, msgid := range valid {
		hash := md5.New()
		hash.Write([]byte(msgid))
		digest := hash.Sum(nil)

		if match.Match(digest) != true {
			t.Fatalf("Msgid(%s) failed", msgid)
		}
	}
	for _, msgid := range invalid {
		hash := md5.New()
		hash.Write([]byte(msgid))
		digest := hash.Sum(nil)

		if match.Match(digest) != false {
			t.Fatalf("Msgid(%s) should fail but didn't", msgid)
		}
	}
}
