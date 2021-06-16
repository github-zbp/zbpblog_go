$(function(){
	 handle_funcs = [tranIp2Domain, tranCode2Markdown, stripNotUseImg, stripPrevTag,stripDataCkeSavedSrc, addDomainToHref, stripCkeTag, stripSpecialStr]
	 var clipboard = new ClipboardJS('#zz', {
		text: function() {
			//获取复制内容
			let cont = handleContent($(".blog-art-cont .cont").html().trim(), handle_funcs);
			let banquan = "本文转载自： " + $(".banquan-info").html();
			return cont+banquan;
		}
	});

	clipboard.on('success', function(e) {
		alert("文章内容已复制到粘贴板中");
	});

	clipboard.on('error', function(e) {
		console.log(e);
	});
	
})

// 对拷贝的内容做修改, f是函数列表，保存用于处理数据的函数
function handleContent(cont, f){
	for(var i in f){
		cont = f[i](cont)
	}
	return cont;
}

// 将链接中的ip改为域名
function tranIp2Domain(cont){
	return cont.replace('81.71.136.86', 'www.zbpblog.com')
}

// 将<code>代码块转为markdown形式
function tranCode2Markdown(cont){
	// 匹配所有的code标签的内容
	cont = cont.replace(/<code[^>]*?>[^]*?<\/code>/g, function(code){
		code = code.replace(/<code[^>].*?>/g, "\n\n```\n")
		code = code.replace(/<\/code>/g, "\n```\n\n")
		code = code.replace(/&nbsp;/g, "")
		var replacement = code.replace(/<\/?\w+[^>]*?>/g, "")
		return replacement
	})
	return cont
}

// 去除不需要的<img>标签
function stripNotUseImg(cont){
	return cont.replace(/<img[^>]*?src="data:[^>]*?>/g, "")
}

// 去除不需要的pre标签
function stripPrevTag(cont){
	cont = cont.replace(/<pre[^>]*?>/g, "")
	cont = cont.replace(/<\/pre>/g, "")
	return cont
}

// 将a标签中所有不带域名的href加上zbpblog.com的域名
function addDomainToHref(cont){
	cont = cont.replace(/<a([^>]*?)href=['"](\/[^>]*?)['"]([^>]*?)>/g, "<a$1href=\"http://www.zbpblog.com$2\"$3>")
	cont = cont.replace(/<img([^>]*?)src=['"](\/[^>]*?)['"]([^>]*?)>/g, "<img$1src=\"http://www.zbpblog.com$2\"$3>")
	return cont
}

// 过滤特殊字符串
function stripSpecialStr(cont){
	return cont.replace("代码段 小部件", "");
}

// 去掉data-cke-saved-src属性
function stripDataCkeSavedSrc(cont){
    return cont.replace(/data-cke-saved-src=".*?"/g, "");
}

// 去掉cke编辑器产生的标签属性（但不去掉标签）
function stripCkeTag(cont){
	return cont.replace(/<(\w+)( [^>]*?cke[^>]*?)>/g, "<$1>")
}