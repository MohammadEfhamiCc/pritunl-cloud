/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as CertificateTypes from '../types/CertificateTypes';
import * as OrganizationTypes from '../types/OrganizationTypes';
import * as CertificateActions from '../actions/CertificateActions';
import * as MiscUtils from '../utils/MiscUtils';
import CertificateDomain from './CertificateDomain';
import PageInput from './PageInput';
import PageSelect from './PageSelect';
import PageInfo from './PageInfo';
import PageTextArea from './PageTextArea';
import PageSave from './PageSave';
import ConfirmButton from './ConfirmButton';
import Help from './Help';
import * as Constants from "../Constants";

interface Props {
	certificate: CertificateTypes.CertificateRo;
	organizations: OrganizationTypes.OrganizationsRo;
}

interface State {
	disabled: boolean;
	changed: boolean;
	message: string;
	certificate: CertificateTypes.Certificate;
	addDomain: string;
}

const css = {
	card: {
		position: 'relative',
		padding: '10px 10px 0 10px',
		marginBottom: '5px',
	} as React.CSSProperties,
	remove: {
		position: 'absolute',
		top: '5px',
		right: '5px',
	} as React.CSSProperties,
	domain: {
		margin: '9px 5px 0 5px',
		height: '20px',
	} as React.CSSProperties,
	itemsLabel: {
		display: 'block',
	} as React.CSSProperties,
	itemsAdd: {
		margin: '8px 0 15px 0',
	} as React.CSSProperties,
	group: {
		flex: 1,
		minWidth: '280px',
		margin: '0 10px',
	} as React.CSSProperties,
	save: {
		paddingBottom: '10px',
	} as React.CSSProperties,
	label: {
		width: '100%',
		maxWidth: '280px',
	} as React.CSSProperties,
	inputGroup: {
		width: '100%',
	} as React.CSSProperties,
};

export default class Certificate extends React.Component<Props, State> {
	constructor(props: any, context: any) {
		super(props, context);
		this.state = {
			disabled: false,
			changed: false,
			message: '',
			certificate: null,
			addDomain: null,
		};
	}

	set(name: string, val: any): void {
		let certificate: any;

		if (this.state.changed) {
			certificate = {
				...this.state.certificate,
			};
		} else {
			certificate = {
				...this.props.certificate,
			};
		}

		certificate[name] = val;

		this.setState({
			...this.state,
			changed: true,
			certificate: certificate,
		});
	}

	onSave = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CertificateActions.commit(this.state.certificate).then((): void => {
			this.setState({
				...this.state,
				message: 'Your changes have been saved',
				changed: false,
				disabled: false,
			});

			setTimeout((): void => {
				if (!this.state.changed) {
					this.setState({
						...this.state,
						message: '',
						changed: false,
						certificate: null,
					});
				}
			}, 3000);
		}).catch((): void => {
			this.setState({
				...this.state,
				message: '',
				disabled: false,
			});
		});
	}

	onDelete = (): void => {
		this.setState({
			...this.state,
			disabled: true,
		});
		CertificateActions.remove(this.props.certificate.id).then((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		}).catch((): void => {
			this.setState({
				...this.state,
				disabled: false,
			});
		});
	}

	onAddDomain = (): void => {
		let cert: CertificateTypes.Certificate;

		if (this.state.changed) {
			cert = {
				...this.state.certificate,
			};
		} else {
			cert = {
				...this.props.certificate,
			};
		}

		let acmeDomains = [
			...cert.acme_domains,
			'',
		];

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: cert,
		});
	}

	onChangeDomain(i: number, state: string): void {
		let cert: CertificateTypes.Certificate;

		if (this.state.changed) {
			cert = {
				...this.state.certificate,
			};
		} else {
			cert = {
				...this.props.certificate,
			};
		}

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains[i] = state;

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			certificate: cert,
		});
	}

	onRemoveDomain(i: number): void {
		let cert: CertificateTypes.Certificate;

		if (this.state.changed) {
			cert = {
				...this.state.certificate,
			};
		} else {
			cert = {
				...this.props.certificate,
			};
		}

		let acmeDomains = [
			...cert.acme_domains,
		];

		acmeDomains.splice(i, 1);

		cert.acme_domains = acmeDomains;

		this.setState({
			...this.state,
			changed: true,
			message: '',
			addDomain: '',
			certificate: cert,
		});
	}

	render(): JSX.Element {
		let cert: CertificateTypes.Certificate = this.state.certificate ||
			this.props.certificate;

		let info: CertificateTypes.Info = this.props.certificate.info || {};

		let organizationsSelect: JSX.Element[] = [];
		organizationsSelect.push(
			<option key="null" value="">
				Node Certificate
			</option>,
		);
		if (this.props.organizations.length) {
			for (let organization of this.props.organizations) {
				organizationsSelect.push(
					<option
						key={organization.id}
						value={organization.id}
					>{organization.name}</option>,
				);
			}
		}

		let domains: JSX.Element[] = [];
		for (let i = 0; i < cert.acme_domains.length; i++) {
			let index = i;

			domains.push(
				<CertificateDomain
					key={index}
					disabled={this.state.disabled || Constants.user}
					domain={cert.acme_domains[index]}
					onChange={(state: string): void => {
						this.onChangeDomain(index, state);
					}}
					onRemove={(): void => {
						this.onRemoveDomain(index);
					}}
				/>,
			);
		}

		return <div
			className="bp3-card"
			style={css.card}
		>
			<div className="layout horizontal wrap">
				<div style={css.group}>
					<div style={css.remove}>
						<ConfirmButton
							safe={true}
							className="bp3-minimal bp3-intent-danger bp3-icon-trash"
							progressClassName="bp3-intent-danger"
							dialogClassName="bp3-intent-danger bp3-icon-delete"
							dialogLabel="Delete Certificate"
							confirmMsg="Permanently delete this certificate"
							disabled={this.state.disabled}
							onConfirm={this.onDelete}
						/>
					</div>
					<PageInput
						label="Name"
						help="Name of certificate"
						type="text"
						placeholder="Name"
						value={cert.name}
						onChange={(val): void => {
							this.set('name', val);
						}}
					/>
					<PageTextArea
						label="Comment"
						help="Certificate comment."
						placeholder="Certificate comment"
						rows={3}
						value={cert.comment}
						onChange={(val: string): void => {
							this.set('comment', val);
						}}
					/>
					<PageTextArea
						readOnly={cert.type !== 'text'}
						label="Private Key"
						help="Certificate private key in PEM format"
						placeholder="Private key"
						rows={6}
						value={cert.key}
						onChange={(val: string): void => {
							this.set('key', val);
						}}
					/>
					<PageTextArea
						readOnly={cert.type !== 'text'}
						label="Certificate Chain"
						help="Certificate followed by certificate chain in PEM format"
						placeholder="Certificate chain"
						rows={6}
						value={cert.certificate}
						onChange={(val: string): void => {
							this.set('certificate', val);
						}}
					/>
					<label
						style={css.itemsLabel}
						hidden={cert.type !== 'lets_encrypt'}
					>
						LetsEncrypt Domains
						<Help
							title="LetsEncrypt Domains"
							content="Enter domain names for the certificate. All domains names must point to a Pritunl Cloud server in the cluster. The servers must also have port 80 publicy open. The port will need to stay open to renew the certificate."
						/>
					</label>
					<div hidden={cert.type !== 'lets_encrypt'}>
						{domains}
					</div>
					<button
						className="bp3-button bp3-intent-success bp3-icon-add"
						disabled={this.state.disabled || Constants.user}
						style={css.itemsAdd}
						hidden={cert.type !== 'lets_encrypt'}
						type="button"
						onClick={this.onAddDomain}
					>
						Add Domain
					</button>
				</div>
				<div style={css.group}>
					<PageInfo
						fields={[
							{
								label: 'ID',
								value: this.props.certificate.id || 'None',
							},
							{
								label: 'Signature Algorithm',
								value: info.signature_alg || 'Unknown',
							},
							{
								label: 'Public Key Algorithm',
								value: info.public_key_alg || 'Unknown',
							},
							{
								label: 'Issuer',
								value: info.issuer || 'Unknown',
							},
							{
								label: 'Issued On',
								value: MiscUtils.formatDate(info.issued_on) || 'Unknown',
							},
							{
								label: 'Expires On',
								value: MiscUtils.formatDate(info.expires_on) || 'Unknown',
							},
							{
								label: 'DNS Names',
								value: info.dns_names || 'Unknown',
							},
						]}
					/>
					<PageSelect
						label="Type"
						disabled={this.state.disabled || Constants.user}
						help="Certificate type, use text to provide a certificate. LetsEncrypt provides free certificates that automatically renew."
						value={cert.type}
						onChange={(val): void => {
							this.set('type', val);
						}}
					>
						<option value="text">Text</option>
						<option value="lets_encrypt">LetsEncrypt</option>
					</PageSelect>
					<PageSelect
						disabled={this.state.disabled}
						hidden={Constants.user}
						label="Organization"
						help="Organization for certificate. Select node to create a certificate for nodes. Load balancers in the same organization will have access to this certificate."
						value={cert.organization}
						onChange={(val): void => {
							this.set('organization', val);
						}}
					>
						{organizationsSelect}
					</PageSelect>
				</div>
			</div>
			<PageSave
				style={css.save}
				hidden={!this.state.certificate}
				message={this.state.message}
				changed={this.state.changed}
				disabled={this.state.disabled}
				light={true}
				onCancel={(): void => {
					this.setState({
						...this.state,
						changed: false,
						certificate: null,
					});
				}}
				onSave={this.onSave}
			/>
		</div>;
	}
}
