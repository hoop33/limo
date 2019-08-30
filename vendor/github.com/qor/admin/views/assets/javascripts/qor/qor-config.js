// init for slideout after show event
$.fn.qorSliderAfterShow = $.fn.qorSliderAfterShow || {};
window.QOR = {
    $formLoading: '<div id="qor-submit-loading" class="clearfix"><div class="mdl-spinner mdl-spinner--single-color mdl-js-spinner is-active"></div></div>'
};

// change Mustache tags from {{}} to [[]]
window.Mustache && (window.Mustache.tags = ['[[', ']]']);

// clear close alert after ajax complete
$(document).ajaxComplete(function(event, xhr, settings) {
    if (settings.type == 'POST' || settings.type == 'PUT') {
        if ($.fn.qorSlideoutBeforeHide) {
            $.fn.qorSlideoutBeforeHide = null;
            window.onbeforeunload = null;
        }
    }
});

// select2 ajax common options
// $.fn.select2 = $.fn.select2 || function(){};
$.fn.select2.ajaxCommonOptions = function(select2Data) {
    let remoteDataPrimaryKey = select2Data.remoteDataPrimaryKey;

    return {
        dataType: 'json',
        headers: getSelect2Header(select2Data),
        cache: true,
        delay: 250,
        data: function(params) {
            return {
                keyword: params.term, // search term
                page: params.page,
                per_page: 20
            };
        },
        processResults: function(data, params) {
            // parse the results into the format expected by Select2
            // since we are using custom formatting functions we do not need to
            // alter the remote JSON data, except to indicate that infinite
            // scrolling can be used
            params.page = params.page || 1;

            var processedData = $.map(data, function(obj) {
                obj.id = obj[remoteDataPrimaryKey] || obj.primaryKey || obj.Id || obj.ID;
                return obj;
            });

            return {
                results: processedData,
                pagination: {
                    more: processedData.length >= 20
                }
            };
        }
    };

};

// select2 ajax common options
// format ajax template data
$.fn.select2.ajaxFormatResult = function(data, tmpl) {
    var result = '';
    if (tmpl.length > 0) {
        result = window.Mustache.render(tmpl.html().replace(/{{(.*?)}}/g, '[[$1]]'), data);
    } else {
        result = data.text || data.Name || data.Title || data.Code || data[Object.keys(data)[0]];
    }

    // if is HTML
    if (/<(.*)(\/>|<\/.+>)/.test(result)) {
        return $(result);
    }
    return result;
};

function getSelect2Header(){
    let data = $("body").data();
    let selectAjaxHeader = data.selectAjaxHeader;
    let getSelect2HeaderFunction = window.getSelect2HeaderFunction;
    let headers = {};

    if(selectAjaxHeader && getSelect2HeaderFunction && $.isFunction(getSelect2HeaderFunction)){
        headers[selectAjaxHeader] = getSelect2HeaderFunction();
        return headers;
    } else {
        return {};
    }
}