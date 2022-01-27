package summary_email_template

func Get() string {
	return `<!DOCTYPE html>
		<html ⚡4email data-css-strict>
		` + head() + `
		` + body() + `
		</html>`
}

func head() string {
	return `<head>
		<meta charset="utf-8" />
		<script async src="https://cdn.ampproject.org/v0.js"></script>
		<script
			async
			custom-template="amp-mustache"
			src="https://cdn.ampproject.org/v0/amp-mustache-0.2.js"
		></script>
		<script
			async
			custom-element="amp-list"
			src="https://cdn.ampproject.org/v0/amp-list-0.1.js"
		></script>
		<script
			async
			custom-element="amp-form"
			src="https://cdn.ampproject.org/v0/amp-form-0.1.js"
		></script>
		<style amp4email-boilerplate>
			body {
				visibility: hidden;
			}
		</style>
		<style amp-custom>
			h1 {
				margin: 1rem;
			}
			.list {
				position: relative;
			}
			.item-wrapper {
				padding: 8px;
			}
			.item {
				display: flex;
				flex-direction: column;
				align-items: space-between;
				padding: 8px;
				border-radius: 16px;
				border: 1px solid black;
			}
			.item div {
				padding-right: 10px;
			}
			.item-head {
				display: flex;
				justify-content: space-between;
				flex-wrap: wrap;
			}
		</style>
	</head>`
}

func body() string {
	return `
	<body>
		<div>
			<h3>Szia Lajos!</h3>
			<p>Számláidat sétáltatod? Hát bazmeg!</p>
		</div>
		<div>
			<amp-list
				src="[[.ApiUrl]]/api/invoices?token=[[.Token]]"
				id="invoiceList"
				binding="no"
				layout="fixed-height"
				single-item
				items="."
				height="10"
			>
				<template type="amp-mustache">
					{{#items}}
					<div class="item-wrapper">
						<div class="item">
							<div class="item-head">
								<div>
									<div>{{accountNumber}}</div>
									<div>{{vendorName}}</div>
								</div>
								<div>{{receivedAt}}</div>
								<div>{{dueDate}}</div>
								<div>{{status}}</div>
								<div>
									<form
										method="post"
										action-xhr="[[.ApiUrl]]/api/invoices?token=[[.Token]]"
										on="submit-success:invoiceList.refresh"
									>
										<input type="hidden" name="invoiceId" value="{{invoiceId}}" />
										<input type="hidden" name="status" value="[[.Accepted]]" />
										<input type="submit" value="Accept" />
										<div submit-success>Accept successful.</div>
										<div submit-error>Accept failed.</div>
									</form>
									<form
										method="post"
										action-xhr="[[.ApiUrl]]/api/invoices?token=[[.Token]]"
										on="submit-success:invoiceList.refresh"
									>
										<input type="hidden" name="invoiceId" value="{{invoiceId}}" />
										<input type="hidden" name="status" value="[[.Rejected]]" />
										<input type="submit" value="Reject" />
										<div submit-success>Reject successful.</div>
										<div submit-error>Reject failed.</div>
									</form>
								</div>
							</div>
							<div class="item-content">
								<div>
									<div>ID: {{invoiceId}}</div>
									<div>Iban: {{iban}}</div>
									<div>InvoiceNumber: {{invoiceNumber}}</div>
									<div>Vendor Email: {{vendorEmail}}</div>
									<div>Net Price: {{netPrice}}</div>
									<div>Gross Price: {{grossPrice}}</div>
									<div>Currency: {{currency}}</div>
									<div>Vat Amount: {{vatAmount}}</div>
									<div>Vat Rate: {{vatRate}}</div>
								</div>
							</div>
						</div>
					</div>
					{{/items}} 
					{{^items}}
					<div>Nincs több pàlinka. Elfogyott.</div>
					<amp-img 
						alt="Pocok logo"
						src="[[.PocokLogo]]"
						width="200px"
						height="200px"
						>
					</amp-img>
					{{/items}}
				</template>
			</amp-list>
		</div>
	</body>`
}
