/* 头部特效 */
// $("#search-select").change(function(){
// 	$("#search-type-input").val($("#search-select option:selected").attr("type"));
// });

$(".search-form").submit(function(){
	if(!$("#search-input").val().trim()){
		return false;
	}
});

// $("#mb-search-select").change(function(){
// 	$("#mb-search-type-input").val($("#mb-search-select option:selected").attr("type"));
// });
$(".mb-search-form").submit(function(){
	if(!$("#mb-search-input").val().trim()){
		return false;
	}
});

$(".header-box-btn").click(function(){
	$(".header-box-modal").hide();
});

$(".mobile-navbar").click(function(){
	$(".header-box-modal").show();
});

$(".mobile-search").click(function(){
	$(".search-modal").show();
});

$(".header-box-ul li").click(function(){
	$(this).find(".header-box-sonul").show();
	console.log(2);
});

$(".mobile-close").click(function(e){
	$(this).parent().hide();
	e.stopPropagation()
	
})

/* 首页加载更多文章特效 */
var index_art_finished=false;
$(".more-blogs").click(function(){
	var page = $(this).attr("data-page");
	if(page==0){
		return false;
	}
	page++;
	var url = "/more_blogs";
	if(index_art_finished){
		$(".more-blogs").text("已加载所有文章");
		//$(".more-blogs").attr("data-page",0);
		return
	}
	$.get(url,{page},function(d){
		if(!d.errno){
			html="";
			if(d.data.data.length){
				$.each(d.data.data,(i,v)=>{
					html+='<li class="item-blog">';
					html+='<div class="blog-img inline-block">';
					html+='<a href="/blog-'+v.id+'.html" target="_blank">';
					html+='<img src="'+v.thumb+'" onerror="this.src=\'/static/home/images/default.png\'" title="'+v.title+'" alt="'+v.title+'"/>';
					
					html+='<a class="cate-tag1" href="/cate-'+v.cate.en_name+'">'+v.cate.zh_name+'</a>';
					html+='</a></div><div class="blog-info inline-block">';
					html+='<h2 class="blog-h2"><a href="/blog-'+v.id+'.html" target="_blank" title="'+v.title+'">'+v.title+'</a></h2>';
					html+='<div>发布时间:<span>'+v.create_time+'</span></div>';
					html+='<p>'+v.description+'</p>';
					html+='<div class="art-tags-box"><ul class="art-tags-ul">';
					if(v.tags && v.tags.length){
						html+='<i class="icon tag-icon"></i>';
						$.each(v.tags,(key,val)=>{
							html+='<li class="float art-tags-li"><a href="/tag-'+val.id+'-'+val.name+'" title="'+val.name+'" target="_blank">'+val.name+'</a></li>';
						});
						
					}
					
					html+='</ul></div></div></li>';
				});
				
				$(".index-ul-blog").append(html);
				
				$(".more-blogs").attr("data-page",page);
			}else{
				index_art_finished=true;
			}
		}
	},"json");
});