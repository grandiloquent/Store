export default class Converter extends HTMLElement {
    constructor() {
        super();
        this.shadowObj = this.attachShadow({mode: 'open'});
        this.render();
    }

    "{1 1}"
    connectedCallback() {
        this.hexDiv = this.shadowObj.getElementById('hexDiv');
        this.hexDiv.addEventListener('input', () => {
            this.onHexInput();
        });
        this.intDiv = this.shadowObj.getElementById('intDiv');
        this.intDiv.addEventListener('input', () => {
            this.onIntInput();
        });
    }

    onHexInput() {
        this.intDiv.value = parseInt(this.hexDiv.value, 16);
    }

    onIntInput() {
        this.hexDiv.value = '0x' + parseInt(this.intDiv.value, 10).toString(16);
    }

    render() {
        this.shadowObj.innerHTML = this.getTemplate();
    }

    getTemplate() {
        return `<input contenteditable="true" class="editable" id="hexDiv">
                <input contenteditable="true" class="editable" id="intDiv">

${this.getStyle()}`;
    }

    getStyle() {
        return `
        <style>
         :host {
        display: flex;
        flex: 1 1 50%;
        flex-direction: row;
        }
        .editable{
           box-shadow: inset 0 1px 2px rgba(27,31,35,.075);
    outline: none;
    border-radius: 3px;
    border: 1px solid #d1d5da;
    background-position: right 8px center;
    background-repeat: no-repeat;
    vertical-align: middle;
    line-height: 20px;
    padding: 6px 8px;
    min-height: 34px;
    font-size: 14px;
    background-color: #fafbfc;
    color: #586069;
    padding-left: 30px;
    width: 320px;
    border-bottom-left-radius: 0;
border-top-left-radius: 0;

}

        </style>
        `
    }
}