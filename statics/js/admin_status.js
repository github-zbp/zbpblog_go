$(function(){
    $(".status").click(function(){
        var cur_status = $(this).attr("data-status")
        var status = cur_status == 0 ? 1 : 0
        var id = $(this).attr("data-id")
        var model = $(this).attr("data-model")
        $.post('/zbpadmin/common/tran_status',{status, id, model}, function(result){
            if(result.errno == 0){
                location.href = location.href
            }else{
                layer.msg('修改状态失败-'+result["errmsg"],{'icon':2});
            }
        }).error(function () {
            layer.msg('系统故障，请联系管理员');
        })
    })
})