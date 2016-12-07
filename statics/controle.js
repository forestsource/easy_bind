"use strict"

$(document).ready(function() {
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
        $("#zone-acl").after(ss.zone);
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
            $("#config-area").html(data);
        }).fail(function(data) { // Error
            $("#config-area").html(data);
        });
    });
});

function read_values() {
    let data = {};
    //zone
    let zones = {};
    let num=0;
    $.each($(".zone-input"), (i, val)=> {
        let zone = {};
        zone.domain = $(val).children(".zone-domain-value").val();
        zone.ip = $(val).children(".zone-ip-value").val();
        if ($(val).children(".zone-cname-check").prop("checked")) {
            zone.cname = $(val).children(".zone-cname-value").val();
        }
        if ($(val).children(".zone-receive-check").prop("checked")) {
            zone.rmail = $(val).children(".zone-receive-mail").val();
        }
        if ($(val).children(".zone-send-check").prop("checked")) {
            zone.smail = $(val).children(".zone-send-mail").val();
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
        acl.domain = $(val).children(".acl-list-name").val();
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
    option.memory= $("#chache-size").val();
    data.option = option;
    console.log(option)//debug
    console.log("==");//debug
    console.log(data);//debug
    return data;
}
var ss = new StaticString()

function StaticString() {
    this.zone = '<div class="zone-input">' +
        '<label class="form-label" for="input-example-1">domain name</label>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="example.com" />' +
        '<label class="form-label" for="input-example-1">ip</label>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="xxx.xxx.xxx.xxx" />' +
        '<input type="checkbox" /> <i class="form-icon"></i>CNAME</br>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="example.com" />' +
        '<input type="checkbox" /> <i class="form-icon"></i>Receiv Mail Server ip</br>' +
        '<input class="form-input " type="text" id="input-example-1" placeholder="mx1.example.com" />' +
        '<input type="checkbox" /> <i class="form-icon"></i>Send Mail server</br>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="xxx.xxx.xxx" />' +
        '<div class="btn btn-sm del-zone-button">-</div>' +
        '</div>'
    this.acl = '<div class="form-group acl-input">' +
        '<label class="form-label" for="input-example-1">List Name</label>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="example.com" />' +
        '<label class="form-label" for="input-example-1">IPs</label>' +
        '<input class="form-input" type="text" id="input-example-1" placeholder="xxx.xxx.xxx.xxx,xxx.xxx.xxx.xxx" />' +
        '<div class="btn btn-sm del-acl-button">-</div>&nbsp' +
        '</div>'
    this.option = ''
}
