$(function(){
	$.getJSON( "data.json", function(data) {
		document.title = "Circa • " + data.folder;

		for (image of data.images) {
			$("#images").append("<img src='/i/"+image+"'></img>");
		}

		$images = $("#images img");
		$section = $("#images");
		$("#images img").click(function() {
			if ($section.hasClass("dimmed")) {
				$section.removeClass("dimmed");
				$images.removeClass("dim");
			} else {
				$section.addClass("dimmed");
				$images.not($(this)).addClass("dim");
			}
		});
	});
});