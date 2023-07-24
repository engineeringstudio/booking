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
    var backend = getDomain();

    layui.use(['laydate', 'jquery'], function () {
        var laydate = layui.laydate;
        var $ = layui.$;

        laydate.render({
            elem: '#date',
            min: 0,
            max: 7,
            showBottom: false,
            ready: function (date) {
                checkQuota($, backend);
            }
        });
    });

    layui.use('form', function () {
        var form = layui.form;

        form.on('submit(form)', function (data) {
            fetch(backend, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json; charset=UTF-8',
                    'credentials': 'include'
                },
                body: JSON.stringify(data.field)
            })
                .then(response => {
                    if (!response.status == 200) {
                        layer.msg("提交失败");
                    }
                    layer.msg("提交成功");
                });
            return false;
        });
    });
}

function getDomain() {
    fetch('config.json', {
        method: 'GET'
    })
        .then(response => {
            if (!response.status == 200) {
                return "api/add";
            }
            return response.json().backend;
        });
}

function checkQuota($ ,backend) {
    var cal = $('#layui-laydate1');
    var dateItems = cal.find('.layui-laydate-content table tbody tr td');
    layui.each(dateItems, function (index, item) {
        if (item.cellIndex == 5 || item.cellIndex == 6) {
            item = $(item);
            item.addClass("laydate-disabled");
            return;
        }

        var date = item.getAttribute('lay-ymd');
        item = $(item);
        fetch(backend + '/check', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json; charset=UTF-8',
                'credentials': 'include'
            },
            body: JSON.stringify({
                date: date
            })
        })
            .then(response => {
                if (response.status == 200) {
                    return response.json();
                }
                return { status: false };
            })
            .then(data => {
                if (data.status == false) {
                    item.addClass("laydate-disabled");
                }
            });
    });
}
