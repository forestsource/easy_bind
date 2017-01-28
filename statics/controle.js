"use strict"

$(document).ready(function() {
    $("#switch-resolver").on("click", function() {
      if($("#switch-resolver").hasClass("active")){
        $("#switch-resolver").removeClass("active");
      }else{
        $("#switch-resolver").addClass("active");
      }
    });
    $("#switch-slave").on("click", function() {
      if($("#switch-slave").hasClass("active")){
        $("#switch-slave").removeClass("active");
      }else{
        $("#switch-slave").addClass("active");
      }
    });
    $("#toggle_zone").on("click", function() {
        $(".tab-item.active").removeClass("active");
        $("section.visible").removeClass("visible");
        $("#toggle_zone").addClass("active");
        $("#zone_area").addClass("visible");
    });
    $("#toggle_acl").on("click", function() {
        $(".tab-item.active").removeClass("active");
        $("section.visible").removeClass("visible");
        $("#toggle_acl").addClass("active");
        $("#acl_area").addClass("visible");
    });
    $("#toggle_key").on("click", function() {
        $(".tab-item.active").removeClass("active");
        $("section.visible").removeClass("visible");
        $("#toggle_key").addClass("active");
        $("#key_area").addClass("visible");
    });
    $("#toggle_option").on("click", function() {
        $(".tab-item.active").removeClass("active");
        $("section.visible").removeClass("visible");
        $("#toggle_option").addClass("active");
        $("#option_area").addClass("visible");
    });
    $(".card").on("click", ".del-zone-button ", function() {
        $(this).parent().remove();
    });
    $(".card").on("click", ".del-acl-button", function() {
        $(this).parent().remove();
    });
    $("#add-zone-button").on("click", function() {
        $(".zone-input").after(ss.zone);
    });
    $("#add-acl-button").on("click", function() {
        $("#acl-head").after(ss.acl);
    });
    $("#make-button").on("click", function() {
        let url = "/make"
        let post_data = {
            ip1: $(".zone-ip-value").val()
        };
        $.ajax({
            type: 'get',
            url: url,
            data: read_values(),
            contentType: 'text/html',
            dataType: 'html',
            scriptCharset: 'utf-8'
        }).done(function(data) { // Success
            console.info(data)
            $("#config-area").val(data);
        }).fail(function(data) { // Error
            $("#config-area").val(data);
        });
    });
});

function read_values() {
    let data = {};
    //toggle
    data.isResolver = $("#switch-resolver").hasClass("active")
    data.isSlave = $("#switch-slave").hasClass("active")
    //zone
    let zones = {};
    let num=0;
    $.each($(".zone-input"), (i, val)=> {
        let zone = {};
        zone.domain = $(val).children(".zone-acl").val();
        zone.ip = $(val).children(".zone-ip-value").val();
        if ($(val).children(".zone-cname-check").prop("checked")) {
            zone.cname = $(val).children(".zone-cname-value").val();
        }
        zone.amail = $(val).children(".zone-admin-value").val();
        if ($(val).children(".is-mail-sever").prop("checked")) {
            zone.isMailServer = "true"
        }
        data["zone"+i] = zone
        ++num;
        console.log(i,zone)//debug
    })
    data.zone_num = num;
    //acl
    let acls = {};
    num=0;
    $.each($(".acl-input"), (i, val) => {
        let acl = {}
        acl.listname = $(val).children(".acl-list-name").val();
        acl.ips = $(val).children(".acl-ip-values").val();
        data["acl"+i] = acl
        console.log(i,acl)//debug
        ++num
    })
    data.acl_num = num;
    //key
    let key = {}
    key.rndc = ($("#rndc-regen").prop("checked")) ? true : false;
    key.tsig = ($("#tsig-regen").prop("checked")) ? true : false;
    key.update = ($("#update-regen").prop("checked")) ? true : false;
    data.key = key;
    console.log(key)//debug
    //option
    let option = {}
    option.ip = $("#maintenance-ip").val();
    option.port = $("#rndc-port").val();
    option.memory= $("#memory-size").val();
    option.isQsyn = ($("#quick-synchronized").prop("checked")) ? true : false;
    option.isMreduce = ($("#memory-reduce").prop("checked")) ? true : false;
    option.isEdns = ($("#edns").prop("checked")) ? true : false;
    option.isResolver = ($("#switch-resolver").hasClass("active")) ? true : false;
    option.isSlave = ($("#switch-slave").hasClass("active")) ? true : false;
    option.forwardIp = $("#forward-ip").val();
    data.option = option;
    console.log(option)//debug
    console.log("==");//debug
    console.log(data);//debug
    return data;
}
var ss = new StaticString()

function StaticString() {
    this.zone = '<div class="zone-input">' +
        '<label class="form-label" for="input-example-1">acl</label>' +
        '<input class="form-input zone-acl" type="text" placeholder="internal" />' +
        '<label class="form-label" for="input-example-1">ip</label>' +
        '<input class="form-input zone-ip-value" type="text" id="input-example-1" placeholder="xxx.xxx.xxx.xxx" />' +
        '<input type="checkbox" /> <i class="form-icon"></i>CNAME</br>' +
        '<input class="form-input zone-cname-value" type="text" id="input-example-1" placeholder="example.com" />' +
        '<label class"form-label">Adminstrator mail adress ip</label></br>' +
        '<input class="form-input zone-admin-value" type="text" id="admin-mail" placeholder="admin@sample.com" />' +
        '<input type="checkbox" class="is-mail-sever"/> <i class="form-icon"></i>Mail Server</br>'+
        '<div class="btn btn-sm del-zone-button">-</div>' +
        '</div>'
    this.acl = '<div class="form-group acl-input">' +
        '<label class="form-label" for="input-example-1">List Name</label>' +
        '<input class="form-input acl-list-name" type="text"  placeholder="example.com" />' +
        '<label class="form-label" for="input-example-1">IPs</label>' +
        '<input class="form-input acl-ip-values" type="text"  placeholder="xxx.xxx.xxx.xxx,xxx.xxx.xxx.xxx" />' +
        '<div class="btn btn-sm del-acl-button">-</div>&nbsp' +
        '</div>'
    this.option = ''
}
