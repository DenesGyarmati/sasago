$(document).ready(function() {
    $('#file_button').on('click', function() {
        $('#file_upload').click();
    });

    $('#file_upload').on('change', function() {
        const fileName = this.files[0] ? this.files[0].name : "No file chosen";
        $('#file_button').text(fileName);
    });

    $('#hamburger').click(function () {
        $(this).toggleClass('active');
        $('#menu').toggleClass('active');
        $('.page').toggleClass('menu-active');
    });

    $('#depth').change(function () {
        const selectedDepth = $(this).val();
        if (selectedDepth === 'Residue' || selectedDepth === 'All') {
            $('#aa-container').removeClass('hidden');
        } else {
            $('#aa-container').addClass('hidden');
        }
    });

    $('#algorithm').change(function () {
        const selectedAlgorithm = $(this).val();
        const parameterInput = $('#parameter');
        if (selectedAlgorithm === 'LR') {
            parameterInput.val(20);
        } else if (selectedAlgorithm === 'SR') {
            parameterInput.val(100);
        }
    });

    $('#advanced').change(function () {
        if ($(this).is(':checked')) {
            $('.form-collector').removeClass('hidden');
        } else {
            $('.form-collector').addClass('hidden');
        }
    });

    $('#use_api').on('click', function() {
        if ($('.form-group.name').hasClass('hidden')) {
            $('.form-group.name').removeClass('hidden');
            $('#file_button').prop('disabled', true);
        } else {
            $('.form-group.name').addClass('hidden');
            $('#file_button').prop('disabled', false);
        }
    });

    $('.accordion-button').click(function () {
        $(this).next('.accordion-content').toggleClass('show');
    });

    $('.residue h3').click(function () {
        $(this).next('.accordion-content').toggleClass('show');
    });

    $('#depth').trigger('change');
});