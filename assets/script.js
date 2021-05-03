$(document).ready(function() {
  document.getElementById('submit').onclick = function() {
    var url = new URL(window.location.href + "generate.png");
    url.searchParams.append("name_beneficiary", $("#name_beneficiary").val());
    url.searchParams.append("iban_beneficiary", $("#iban_beneficiary").val());
    url.searchParams.append("amount", $("#amount").val());
    url.searchParams.append("remittance", $("#remittance").val());
    if ($("#structured").prop("checked")) {
      url.searchParams.append("structured", "1");
    }

    $.ajax({
      type: "GET",
      url: url,
      success: function() {
        $("#reply").text("");
        $("#codeImage").attr("src", url);
      },
    }).fail(function(XMLHttpRequest, textStatus, errorThrown) {
      $("#reply").text(XMLHttpRequest.responseJSON.message);
    });
  }
});
