jQuery(function ($) {
	$("body.translate form#form-translate input").change(function () {
		var input = $(this);
		var name = input.attr('name');
		var original = $("#"+name+"-original").val();
		var partOf = $("#"+name+"-partof").val();
		var translation = input.val();
		$.get('/api/translate', {
			language: $("#current-language").val(),
			original: original,
			partOf: partOf,
			translation: translation
		}, function () {
			input.closest("tr.entry").removeClass("untranslated");
			$("#saved-notice").stop(true).show().fadeTo(0, 1.0).fadeOut(2500);
		});
	}).focus(function () {
		$(this).closest('td.translation-parts').find('.my-translation-arrow-score').addClass('focus');
		$(this).closest('tr.entry').addClass('focus');
		changeLookup();
	}).blur(function () {
		$(this).closest('td.translation-parts').find('.my-translation-arrow-score').removeClass('focus');
		$(this).closest('tr.entry').removeClass('focus');
	});

	$("body.translate a.vote").click(function () {
		var a = $(this);
		var up = a.is('.vote-up');
		var active = a.is('.active');

		var original = a.closest('tr.entry').find('input.entry-original').first().val();
		var partOf = a.closest('tr.entry').find('input.entry-partof').first().val();
		var translation = a.closest('.other-translation').find('label.part').first().text();

		$.get('/api/vote', {
			language: $("#current-language").val(),
			original: original,
			partOf: partOf,
			translation: translation,
			up: (up && !active),
			down: (!up && !active)
		});

		if (up) {
			a.closest("td").find("a.vote-up").removeClass('active btn-success');
		}

		var inverse = a.closest('.btn-group').find('a.vote-'+(up ? 'down' : 'up'));
		inverse.removeClass('active btn-success btn-danger');
		
		if (active) {
			a.removeClass('active btn-success btn-danger');
		} else {
			a.addClass('active btn-'+(up ? 'success' : 'danger'));
		}
	});

	var pageOptions = $("form.page-options");
	pageOptions.find("input, select").change(function () {
		pageOptions.submit();
	});

	$("a.api").click(function () {
		var a = $(this);
		var href = a.attr('href');
		$.get(href, a.data(), function () {
			if (a.is(".reload")) {
				location.reload(true);
			}
		});
		return false;
	});

	$("a.reveal-my-translation").click(function () {
		$(this).closest("tr").addClass("with-translation").find("p.my-translation input").first().focus();
	});

	var lastLookup = '';
	function changeLookup() {
		var lookup = $("#lookup").val();
		if (lookup == "") {
			lookup = $("tr.entry.focus").data("entry-original-text");
			$("#lookup").attr('placeholder', lookup);
		}
		if (lookup == lastLookup)
			return;
		lastLookup = lookup;
		if (lookup == "" || typeof lookup === 'undefined')
			return;

		// ajax lookup
		lookup = lookup.replace("|", " ");
		$.get('/api/lookup', {
			language: $("#current-language").val(),
			lookup: lookup,
		}, function (data) {
			$("#lookup-results").html(data);
			layoutFooterSpacer();
		});
	}

	$("#close-translation-hint").click(function (e) {
		e.preventDefault();
		$("#lookup-results").html("");
		$("#lookup").val("");
		layoutFooterSpacer();
	});

	function layoutFooterSpacer() {
		var height = $("#translation-hint").outerHeight();
		$("#footer-spacer").css({height: height+"px"});
	}

	$("#lookup").change(changeLookup).blur(changeLookup);

	$("#lookup-form").submit(function (e) {
		e.preventDefault();
		changeLookup();
	});

	layoutFooterSpacer();
});