$(function(){
    $(".del").click(function(){
        if(confirm("是否确认删除?")){
            $.get($(this).attr("href"),{},function(d){
                if(!d.errno){
                    location.href=location.href;
                }
            },"json");
        }

        return false;
    });


})

function bindDelMany(model){
    $(".del_many").click(function(){
        var checked_box = $(".del-checkbox:checked");
        var ids = [];

        checked_box.each(function(){
            ids.push($(this).attr('data-id'));
        });
console.log(ids);
        if(ids.length==0 || !confirm("确认批量删除?")){
            return false;
        }


        $.post("/zbpadmin/common/del",{ids, model},function(d){
            if(!d.errno){
                location.href=location.href;
            }
        },"json");
    })
}