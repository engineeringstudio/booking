/**
 * Copyright (C) 2023 CharlieYu4994
 * 
 * This file is part of booking.
 * 
 * booking is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * booking is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with booking.  If not, see <http://www.gnu.org/licenses/>.
 */

function layuiInit() {
    layui.use('laydate', function () {
        var laydate = layui.laydate;

        laydate.render({
            elem: '#date'
        });
    });

    layui.use('form', function () {
        var form = layui.form;

        form.on('submit(form)', function (data) {
            var request = new XMLHttpRequest();
            var backend = getDomain();

            request.withCredentials = true;
            request.open("POST", backend);
            request.setRequestHeader('Content-Type', 'application/json; charset=UTF-8')
            request.send(JSON.stringify(data.field));
            request.addEventListener("load", function () {
                if (request.status == 200) {
                    layer.msg("提交成功");
                } else {
                    layer.msg("提交失败");
                }
            });
            return false;
        });
    });
}


function getDomain() {
    var url = "config.json"
    var request = new XMLHttpRequest();
    request.open("get", url, false);
    request.send(null);
    if (request.readyState == 4) {
        if (request.status == 200) {
            return JSON.parse(request.responseText)["backend"];
        } else {
            return "api/add";
        }
    }
}
